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

package tasks_test

import (
	"x-patrol/tasks"
	"x-patrol/models"
	"x-patrol/vars"
	"x-patrol/util/index"

	"testing"
	"encoding/json"
)

func TestGenerateSearcher(t *testing.T) {
	reposCfg := make([]models.RepoConfig, 0)
	repoConfig := models.RepoConfig{Name: "xsec-traffic", Url: "https://github.com/netxfly/xsec-traffic",
		PollInterval: 30, Vcs: "git", UrlPattern: models.UrlPattern{BaseUrl: vars.DefaultBaseUrl, Anchor: vars.DefaultAnchor},
		AutoPullUpdate: true, ExcludeDotFiles: true,
	}

	reposCfg = append(reposCfg, repoConfig)

	t.Log(tasks.GenerateSearcher(reposCfg))
}

func TestSearchRepos(t *testing.T) {
	reposCfg := make([]models.RepoConfig, 0)
	repoConfig := models.RepoConfig{Name: "xsec-traffic", Url: "https://github.com/netxfly/xsec-traffic",
		PollInterval: 30, Vcs: "git", UrlPattern: models.UrlPattern{BaseUrl: vars.DefaultBaseUrl, Anchor: vars.DefaultAnchor},
		AutoPullUpdate: true, ExcludeDotFiles: true,
	}

	var filesOpened int
	var durationMs int
	reposCfg = append(reposCfg, repoConfig)
	searchers, errors, hasError, err := tasks.GenerateSearcher(reposCfg)
	t.Log(searchers, errors, hasError, err)

	repos := make([]string, 0)
	repos = append(repos, "xsec-traffic")
	rule := models.Rules{Part: "keyword", Type: "regex", Pattern: "password", Caption: "Contains word: password",
		Description: "Contains word: password"}
	opts := index.SearchOptions{IgnoreCase: true, LinesOfContext: tasks.DefaultLinesOfContext, FileRegexp: ""}
	respSearch, err := tasks.SearchRepos(rule, &opts, repos, searchers, &filesOpened, &durationMs)
	respJson, err := json.Marshal(respSearch)
	t.Log(string(respJson), err)
}
