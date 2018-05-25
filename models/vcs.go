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

package models

import (
	"x-patrol/vars"

	"time"
)

type (
	UrlPattern struct {
		Id      int64
		BaseUrl string `json:"base-url"`
		Anchor  string `json:"anchor"`
		Vcs     string `json:"vcs"`
	}

	VcsConfig struct {
		Id       int64
		Username string `json:"username"`
		Password string `json:"password"`
	}

	RepoConfig struct {
		Id              int64
		Name            string        `json:"name"`
		Url             string        `json:"url"`
		PollInterval    time.Duration `json:"poll_interval"`
		Vcs             string        `json:"vcs"`
		UrlPattern      UrlPattern    `json:"url_pattern"`
		VcsConfig       VcsConfig     `json:"vcs_config"`
		AutoPullUpdate  bool          `json:"auto_pull_update"`
		ExcludeDotFiles bool          `json:"exclude_dot_files"`
		Status          int           `json:"status"`
	}
)

func NewUrlParttern(baseUrl, anchor, vcs string) (*UrlPattern) {
	return &UrlPattern{BaseUrl: baseUrl, Anchor: anchor, Vcs: vcs}
}

func (u *UrlPattern) Exist() (bool, error) {
	return Engine.Get(u)
}

func (u *UrlPattern) Insert() (int64, error) {
	return Engine.Insert(u)
}

func ListUrlPattern() ([]UrlPattern, error) {
	UrlPatterns := make([]UrlPattern, 0)
	err := Engine.Find(&UrlPatterns)
	return UrlPatterns, err
}

func InitUrlPattern() () {
	urlPatterns := make([]*UrlPattern, 0)

	gitPattern := NewUrlParttern(vars.DefaultBaseUrl, vars.DefaultAnchor, "github")
	gitlabPattern := NewUrlParttern("{url}/master/{path}{anchor}", vars.DefaultAnchor, "gitlab")
	svnPattern := NewUrlParttern("{url}/{path}{anchor}", "", "svn")
	// localPatten := NewUrlParttern(vars.DefaultBaseUrl, vars.DefaultAnchor, "git")
	bitbucketPatten := NewUrlParttern("{url}/src/master/{path}{anchor}", "#{filename}-{line}", "bitbucket")

	urlPatterns = append(urlPatterns, gitPattern)
	urlPatterns = append(urlPatterns, gitlabPattern)
	urlPatterns = append(urlPatterns, svnPattern)
	// urlPatterns = append(urlPatterns, localPatten)
	urlPatterns = append(urlPatterns, bitbucketPatten)

	for _, u := range urlPatterns {
		has, err := u.Exist()
		if err == nil && !has {
			u.Insert()
		}
	}
}

func GetUrlPattern(vcs string) (bool, error, *UrlPattern) {
	u := new(UrlPattern)
	has, err := Engine.Table("url_pattern").Where("vcs=?", vcs).Get(u)
	return has, err, u
}

func NewRepoConfig(name string,
	url string,
	interval time.Duration,
	vcs string,
	urlPat UrlPattern,
	vcsConf VcsConfig,
	isPull bool,
	isExclude bool) (*RepoConfig) {
	return &RepoConfig{Name: name, Url: url, PollInterval: interval, Vcs: vcs, UrlPattern: urlPat,
		VcsConfig: vcsConf, AutoPullUpdate: isPull, ExcludeDotFiles: isExclude, Status: 1}
}

func (r *RepoConfig) Insert() (int64, error) {
	return Engine.Insert(r)
}

func (r *RepoConfig) Exist() (bool, error) {
	has := false
	rs := make([]RepoConfig, 0)
	err := Engine.Where("url=?", r.Url).Find(&rs)
	if err == nil && len(rs) > 0 {
		has = true
	}
	return has, err
}

func ListRepoConfig() ([]RepoConfig, error) {
	reposConfig := make([]RepoConfig, 0)
	err := Engine.Find(&reposConfig)
	return reposConfig, err
}

func ListRepoConfigPage(page int) ([]RepoConfig, int, error) {
	reposConf := make([]RepoConfig, 0)
	totalPages, err := Engine.Table("repo_config").Count()
	var pages int

	if int(totalPages)%vars.PAGE_SIZE == 0 {
		pages = int(totalPages) / vars.PAGE_SIZE
	} else {
		pages = int(totalPages)/vars.PAGE_SIZE + 1
	}

	if page >= pages {
		page = pages
	}

	if page < 1 {
		page = 1
	}

	err = Engine.Limit(vars.PAGE_SIZE, (page-1)*vars.PAGE_SIZE).Find(&reposConf)
	return reposConf, pages, err
}

func ListValidRepoConfig() ([]RepoConfig, error) {
	reposConfig := make([]RepoConfig, 0)
	err := Engine.Where("status=1").Find(&reposConfig)
	return reposConfig, err
}

func EnableRepoConfById(id int64) (error) {
	repoConf := new(RepoConfig)
	has, err := Engine.ID(id).Get(repoConf)
	if err == nil && has {
		repoConf.Status = 1
		_, err = Engine.ID(id).Cols("status").Update(repoConf)
	}
	return err
}

func DisableRepoConfById(id int64) (error) {
	repoConf := new(RepoConfig)
	has, err := Engine.ID(id).Get(repoConf)
	if err == nil && has {
		repoConf.Status = 0
		_, err = Engine.ID(id).Cols("status").Update(repoConf)
	}
	return err
}

func DeleteAllReposConf() (error) {
	sqlCmd := "delete from repo_conf"
	_, err := Engine.Exec(sqlCmd)
	return err
}

func DeleteRepoConfById(id int64) (err error) {
	repoConf := new(RepoConfig)
	_, err = Engine.Id(id).Delete(repoConf)
	return err
}
