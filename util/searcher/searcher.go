/*

Copyright (c) 2018 sec.xiaomi.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package searcher

import (
	"x-patrol/util/index"
	"x-patrol/models"
	"x-patrol/logger"
	"x-patrol/vars"
	"x-patrol/util/vcs"

	"sync"
	"time"
	"path/filepath"
	"io/ioutil"
	"runtime"
	"os"
	"fmt"
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
)

type empty struct{}
type limiter chan bool

type Searcher struct {
	idx        *index.Index
	lck        sync.RWMutex
	RepoConfig *models.RepoConfig

	// The channel is used to request updates from the API and
	// to signal that it is ok for searchers to begin polling.
	// It has a buffer size of 1 to allow at most one pending
	// update at a time.
	updateCh chan time.Time

	shutdownRequested bool
	shutdownCh        chan empty
	doneCh            chan empty
}

/**
 * Holds a set of IndexRefs that were found in the dbpath at startup,
 * these indexes can be 'claimed' and re-used by newly created searchers.
 */
type foundRefs struct {
	refs    []*index.IndexRef
	claimed map[*index.IndexRef]bool
}

func makeLimiter(n int) limiter {
	return limiter(make(chan bool, n))
}

func (l limiter) Acquire() {
	l <- true
}

func (l limiter) Release() {
	<-l
}

/**
 * Find an Index ref for the repo url and rev, returns nil if no such
 * ref exists.
 */
func (r *foundRefs) find(url, rev string) *index.IndexRef {
	for _, ref := range r.refs {
		if ref.Url == url && ref.Rev == rev {
			return ref
		}
	}
	return nil
}

/**
 * Claim a ref for reuse. This ensures they ref will not be garbage
 * collected at the end of startup.
 */
func (r *foundRefs) claim(ref *index.IndexRef) {
	r.claimed[ref] = true
}

/**
 * Delete the directorires associated with all IndexRefs that were
 * found in the dbpath but were not claimed during startup.
 */
func (r *foundRefs) removeUnclaimed() error {
	for _, ref := range r.refs {
		if r.claimed[ref] {
			continue
		}

		if err := ref.Remove(); err != nil {
			return err
		}
	}
	return nil
}

// Perform atomic swap of index in the searcher so that the new
// index is made "live".
func (s *Searcher) swapIndexes(idx *index.Index) error {
	s.lck.Lock()
	defer s.lck.Unlock()

	oldIdx := s.idx
	s.idx = idx

	return oldIdx.Destroy()
}

// Perform a basic search on the current index using the supplied pattern
// and the options.
//
// TODO(knorton): pat should really just be a part of SearchOptions
func (s *Searcher) Search(pat string, opt *index.SearchOptions) (*index.SearchResponse, error) {
	s.lck.RLock()
	defer s.lck.RUnlock()
	return s.idx.Search(pat, opt)
}

// Get the excluded files as a JSON string. This is only used for returning
// the data directly to clients (thus JSON).
func (s *Searcher) GetExcludedFiles() string {
	path := filepath.Join(s.idx.GetDir(), "excluded_files.json")
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Log.Errorf("Couldn't read excluded_files.json %v\n", err)
	}
	return string(dat)
}

// Triggers an immediate poll of the repository.
func (s *Searcher) Update() bool {
	if !s.RepoConfig.AutoPullUpdate {
		return false
	}

	// schedule an update if one is not already scheduled
	select {
	case s.updateCh <- time.Now():
	default:
		// don't wait to enqueue another update
	}

	return true
}

// Shut down the searcher cleanly, waiting for any indexing operations to complete.
func (s *Searcher) Stop() {
	select {
	case s.shutdownCh <- empty{}:
		s.shutdownRequested = true
	default:
	}
}

// Blocks until the searcher's associated goroutine is stopped.
func (s *Searcher) Wait() {
	<-s.doneCh
}

func (s *Searcher) completeShutdown() {
	close(s.doneCh)
}

// Wait for either the delay period to expire or an update request to
// arrive. Note that an empty delay will result in an infinite timeout.
func (s *Searcher) waitForUpdate(delay time.Duration) {
	var tch <-chan time.Time
	if delay.Nanoseconds() > 0 {
		tch = time.After(delay)
	}

	// wait for a timeout, the update channel signal, or a shutdown request
	select {
	case <-s.updateCh:
	case <-tch:
	case <-s.shutdownCh:
	}
}

// Signal the searcher that it is ok to begin polling the repository.
func (s *Searcher) begin() {
	s.updateCh <- time.Now()
}

// Generate a new index directory in the dbpath. The names are based
// on pseudo-randomness with a time-based seed.
func nextIndexDir(dbpath string) string {
	r := uint64(rand.Uint32())<<32 | uint64(rand.Uint32())
	return filepath.Join(dbpath, fmt.Sprintf("idx-%08x", r))
}

// Read the refs associated with each of the index dirs
// in the given dbpath.
func findExistingRefs(dbpath string) (*foundRefs, error) {
	dirs, err := filepath.Glob(filepath.Join(dbpath, "idx-*"))
	if err != nil {
		return nil, err
	}

	var refs []*index.IndexRef
	for _, dir := range dirs {
		r, _ := index.Read(dir)
		refs = append(refs, r)
	}

	return &foundRefs{
		refs:    refs,
		claimed: map[*index.IndexRef]bool{},
	}, nil
}

// Open an index at the given path. If the idxDir is already present, it will
// simply open and use that index. If, however, the idxDir does not exist a new
// one will be built.
func buildAndOpenIndex(
	opt *index.IndexOptions,
/*dbpath,*/
	vcsDir,
	idxDir,
	url,
	rev string) (*index.Index, error) {
	if _, err := os.Stat(idxDir); err != nil {
		r, err := index.Build(opt, idxDir, vcsDir, url, rev)
		if err != nil {
			return nil, err
		}

		return r.Open()
	}

	return index.Open(idxDir)
}

// Simply prints out statistics about the heap. When hound rebuilds a new
// index it will expand the heap with a decent amount of garbage. This is
// helpful to ensure the heap growth looks sane.
func reportOnMemory() {
	var ms runtime.MemStats

	// Print out interesting heap info.
	runtime.ReadMemStats(&ms)
	fmt.Printf("HeapInUse = %0.2f\n", float64(ms.HeapInuse)/1e6)
	fmt.Printf("HeapIdle  = %0.2f\n", float64(ms.HeapIdle)/1e6)
}

// Utility function for producing a hex encoded sha1 hash for a string.
func hashFor(name string) string {
	h := sha1.New()
	h.Write([]byte(name))
	return hex.EncodeToString(h.Sum(nil))
}

// Create a normalized name for the vcs directory of this repo.
func vcsDirFor(repoCfg *models.RepoConfig) string {
	return fmt.Sprintf("vcs-%s", hashFor(repoCfg.Url))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Make a searcher for each repo in the Config. This function kind of has a notion
// of partial errors. First, if the error returned is non-nil then a fatal error has
// occurred and no other return values are valid. If an error occurs that is specific
// to a particular searcher, that searcher will not be present in the searcher map and
// will have an error entry in the error map.
func MakeAll(reposCfg []models.RepoConfig) (map[string]*Searcher, map[string]error, error) {
	errs := map[string]error{}
	searchers := map[string]*Searcher{}

	refs, err := findExistingRefs(vars.REPO_PATH)
	if err != nil {
		return nil, nil, err
	}

	lim := makeLimiter(vars.MAX_INDEXERS)

	for _, repoCfg := range reposCfg {
		s, err := newSearcher(vars.REPO_PATH, repoCfg.Name, &repoCfg, refs, lim)
		if err != nil {
			logger.Log.Errorf("open file: %v error: %v", repoCfg.Name, err)
			errs[repoCfg.Name] = err
			continue
		}

		searchers[repoCfg.Name] = s
	}

	if err := refs.removeUnclaimed(); err != nil {
		return nil, nil, err
	}

	// after all the repos are in good shape, we start their polling
	for _, s := range searchers {
		s.begin()
	}

	return searchers, errs, nil
}

// Creates a new Searcher that is available for searches as soon as this returns.
// This will pull or clone the target repo and start watching the repo for changes.
func New(dbpath, name string, repoCfg *models.RepoConfig) (*Searcher, error) {
	s, err := newSearcher(dbpath, name, repoCfg, &foundRefs{}, makeLimiter(1))
	if err != nil {
		return nil, err
	}

	s.begin()

	return s, nil
}

// Update the vcs and reindex the given repo.
func updateAndReindex(
	s *Searcher,
	dbpath,
	vcsDir,
	name,
	rev string,
	wd *vcs.WorkDir,
	opt *index.IndexOptions,
	lim limiter) (string, bool) {

	// acquire a token from the rate limiter
	lim.Acquire()
	defer lim.Release()

	repoConfig := s.RepoConfig
	newRev, err := wd.PullOrClone(vcsDir, repoConfig.Url)

	if err != nil {
		logger.Log.Printf("vcs pull error (%s - %s): %s", name, repoConfig.Url, err)
		return rev, false
	}

	if newRev == rev {
		return rev, false
	}

	//logger.Log.Printf("Rebuilding %s for %s", name, newRev)
	idx, err := buildAndOpenIndex(
		opt,
		vcsDir,
		nextIndexDir(dbpath),
		repoConfig.Url,
		newRev)
	if err != nil {
		logger.Log.Errorf("failed index build (%s): %s", name, err)
		return rev, false
	}

	if err := s.swapIndexes(idx); err != nil {
		logger.Log.Errorf("failed index swap (%s): %s", name, err)
		if err := idx.Destroy(); err != nil {
			logger.Log.Errorf("failed to destroy index (%s): %s\n", name, err)
		}
		return rev, false
	}

	return newRev, true
}

// Creates a new Searcher that is capable of re-claiming an existing index directory
// from a set of existing manifests.
func newSearcher(
	dbPath, name string,
	repoConfig *models.RepoConfig,
	refs *foundRefs,
	lim limiter) (*Searcher, error) {

	vcsDir := filepath.Join(dbPath, vcsDirFor(repoConfig))

	wd, err := vcs.New(repoConfig.Vcs, nil)
	if err != nil {
		return nil, err
	}

	opt := &index.IndexOptions{
		ExcludeDotFiles: repoConfig.ExcludeDotFiles,
		SpecialFiles:    wd.SpecialFiles(),
	}

	rev, err := wd.PullOrClone(vcsDir, repoConfig.Url)
	if err != nil {
		return nil, err
	}

	var idxDir string
	ref := refs.find(repoConfig.Url, rev)
	if ref == nil {
		idxDir = nextIndexDir(dbPath)
	} else {
		idxDir = ref.Dir()
		refs.claim(ref)
	}

	idx, err := buildAndOpenIndex(
		opt,
		vcsDir,
		idxDir,
		repoConfig.Url,
		rev)
	if err != nil {
		return nil, err
	}

	s := &Searcher{
		idx:        idx,
		updateCh:   make(chan time.Time, 1),
		RepoConfig: repoConfig,
		doneCh:     make(chan empty),
		shutdownCh: make(chan empty, 1),
	}

	go func() {

		// each searcher's poller is held until begin is called.
		<-s.updateCh

		// if all forms of updating are turned off, we're done here.
		if !repoConfig.AutoPullUpdate {
			s.completeShutdown()
			return
		}

		var delay time.Duration
		if repoConfig.AutoPullUpdate {
			delay = time.Duration(repoConfig.PollInterval) * time.Second
		}

		for {
			// Wait for a signal to proceed
			s.waitForUpdate(delay)

			if s.shutdownRequested {
				s.completeShutdown()
				return
			}

			// attempt to update and reindex this searcher
			newRev, ok := updateAndReindex(s, dbPath, vcsDir, name, rev, wd, opt, lim)
			if !ok {
				continue
			}

			rev = newRev

			// This is just a good time to GC since we know there will be a
			// whole set of dead posting lists on the heap. Ensuring these
			// go away quickly helps to prevent the heap from expanding
			// uncessarily.
			runtime.GC()

			reportOnMemory()
		}
	}()

	return s, nil
}
