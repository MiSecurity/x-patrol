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

package tasks

import (
	"x-patrol/util/index"
	"x-patrol/util/searcher"
	"x-patrol/models"
	"x-patrol/vars"
	"x-patrol/logger"
	"x-patrol/util/lib"

	"time"
	"os"
	"sync"
	"strings"
)

const (
	DefaultLinesOfContext uint = 4
	MaxLinesOfContext     uint = 20
)

type Stats struct {
	FilesOpened int
	Duration    int
}

type SearchResponse struct {
	repo string
	res  *index.SearchResponse
	err  error
}

func init() {
	vars.Exts = make(map[string]bool)
	vars.Exts["js"] = true
	vars.Exts["css"] = true
	vars.Exts["html"] = true
	vars.Exts["htm"] = true
}

func GenerateSearcher(reposConfig []models.RepoConfig) (map[string]*searcher.Searcher, map[string]bool, bool, error) {
	errRepos := make(map[string]bool)
	hasError := false

	// Ensure we have a repos path
	if _, err := os.Stat(vars.REPO_PATH); err != nil {
		if err := os.MkdirAll(vars.REPO_PATH, os.ModePerm); err != nil {
			return nil, errRepos, hasError, err
		}
	}

	searchers, errs, err := searcher.MakeAll(reposConfig)
	if err == nil {
		if len(errs) > 0 {
			// NOTE: This mutates the original config so the repos
			// are not even seen by other code paths.

			for name := range errs {
				errRepos[name] = true
			}
			hasError = true
		}
	}

	return searchers, errRepos, hasError, nil
}

/* search repos, return a  map[string]*index.SearchResponse map */
func SearchRepos(
	rule models.Rules,
	opts *index.SearchOptions,
	repos []string,
	idx map[string]*searcher.Searcher,
	filesOpened *int,
	duration *int,
) (map[string]*index.SearchResponse, error) {
	query := rule.Pattern
	startedAt := time.Now()
	num := len(repos)

	// use a buffered channel to avoid routine leaks on errs.
	ch := make(chan *SearchResponse, num)
	for _, repo := range repos {
		go func(repo string) {
			logger.Log.Infof("check repo: %v with rule: %v, filename: %v", repo, query, opts.FileRegexp)
			fms, err := idx[repo].Search(query, opts)
			ch <- &SearchResponse{repo, fms, err}
		}(repo)
	}

	res := map[string]*index.SearchResponse{}
	for i := 0; i < num; i++ {

		r := <-ch
		if r == nil {
			continue
		}

		if r.res == nil {
			continue
		}

		r.res.RuleId = rule.Id
		r.res.RuleCaption = rule.Caption
		r.res.RulePattern = rule.Pattern

		if r.err != nil {
			return nil, r.err
		}

		if r.res.Matches == nil {
			continue
		}
		res[r.repo] = r.res
		*filesOpened += r.res.FilesOpened
	}

	*duration = int(time.Now().Sub(startedAt).Seconds() * 1000)

	return res, nil
}

func DoSearch(reposConfig []models.RepoConfig, rules models.Rules) (map[string]*index.SearchResponse, models.Rules, error) {
	searchers, errors, _, err := GenerateSearcher(reposConfig)
	respSearch := make(map[string]*index.SearchResponse)
	if err == nil {
		repos := make([]string, 0)
		for _, repoCfg := range reposConfig {
			repo := repoCfg.Name
			if !errors[repo] {
				repos = append(repos, repoCfg.Name)
			}
		}

		opts := index.SearchOptions{IgnoreCase: true, LinesOfContext: DefaultLinesOfContext}
		if strings.ToLower(rules.Part) == "keyword" {
			// search keyword from all files
			opts.FileRegexp = ""
		} else {
			// when rules.Part in ("filename", "path", "extension"), only search filename, and set rules.Pattern = "\\."
			opts.FileRegexp = rules.Pattern
			rules.Pattern = "\\."
		}

		var filesOpened int
		var durationMs int

		respSearch, err = SearchRepos(rules, &opts, repos, searchers, &filesOpened, &durationMs)
	}
	return respSearch, rules, err
}

// 分割任务为map形式，key为批次，value为一批models.RepoConfig
func SegmentationTask(reposConfig []models.RepoConfig) (map[int][]models.RepoConfig) {
	tasks := make(map[int][]models.RepoConfig)
	totalRepos := len(reposConfig)
	scanBatch := totalRepos / vars.MAX_Concurrency_REPOS

	for i := 0; i < scanBatch; i++ {
		curTask := reposConfig[vars.MAX_Concurrency_REPOS*i : vars.MAX_Concurrency_REPOS*(i+1)]
		tasks[i] = curTask
	}

	if totalRepos%vars.MAX_Concurrency_REPOS > 0 {
		n := len(tasks)
		tasks[n] = reposConfig[vars.MAX_Concurrency_REPOS*scanBatch : totalRepos]
	}
	return tasks
}

// 按批次分发、执行任务
func DistributionTask(tasksMap map[int][]models.RepoConfig, rules []models.Rules) {
	for _, rule := range rules {
		for _, reposConf := range tasksMap {
			Run(reposConf, rule)

		}
	}
}

func Run(reposConfig []models.RepoConfig, rule models.Rules) {
	var wg sync.WaitGroup
	wg.Add(len(reposConfig))
	for _, rConfig := range reposConfig {
		reposCfg := make([]models.RepoConfig, 0)
		reposCfg = append(reposCfg, rConfig)

		go func(config []models.RepoConfig, rule models.Rules) {
			defer wg.Done()
			SaveSearchResult(DoSearch(reposCfg, rule))
		}(reposCfg, rule)
	}
	// wg.Wait()
	waitTimeout(&wg, vars.TIME_OUT*time.Second)
}

func SaveSearchResult(responses map[string]*index.SearchResponse, rule models.Rules, err error, ) {
	if err == nil {
		for repo, resp := range responses {
			revision := resp.Revision
			for _, fileMatches := range resp.Matches {

				filename := fileMatches.Filename
				ext := GetExt(filename)
				if vars.Exts[ext] {
					continue
				}

				for _, matches := range fileMatches.Matches {
					hash := lib.MakeHash(repo, revision, filename, matches.Line)
					// logger.Log.Infof("repo:%v, revision:%v, filename: %v, matches Line: %v", repo,
					//	revision, filename, matches.Line)
					result := models.NewSearchResult(matches, repo, filename, revision, hash, rule)
					has, err := result.Exist()
					if err == nil && ! has {
						result.Insert()
					}
				}
			}
		}
	}
}

func ScheduleTasks(duration time.Duration) () {
	for {
		// insert repos from inputInfo
		// githubsearch.InsertAllRepos()
		// insert all enable repos to repos config table
		// models.InsertReposConfig()

		rules, err := models.GetLocalRules()
		if err == nil {
			reposConfig, err := models.ListValidRepoConfig()
			if err == nil {
				mapTasks := SegmentationTask(reposConfig)
				DistributionTask(mapTasks, rules)
			}
		}

		logger.Log.Infof("Complete the scan local repos, start to sleep %v seconds", duration*time.Second)

		time.Sleep(duration * time.Second)
	}
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

func GetExt(filename string) (ext string) {
	exts := strings.Split(filename, ".")

	if len(exts) > 1 {
		ext = exts[len(exts)-1]
	}
	return ext
}
