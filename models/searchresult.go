/*

Copyright (c) 2018 sec.lu

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

	"github.com/google/go-github/github"

	"time"
)

type Match struct {
	Id      int64
	Text    *string `json:"text,omitempty" xorm:"LONGBLOB"`
	Indices []int   `json:"indices,omitempty" xorm:"json"`
}

// TextMatch represents a text match for a SearchResult
type TextMatch struct {
	Id         int64
	ObjectURL  *string `json:"object_url,omitempty"`
	ObjectType *string `json:"object_type,omitempty"`
	Property   *string `json:"property,omitempty"`
	Fragment   *string `json:"fragment,omitempty"`
	Matches    []Match `xorm:"LONGBLOB"`
}

// CodeResult represents a single search result.
type CodeResult struct {
	Id          int64
	Name        *string            `json:"name,omitempty"`
	Path        *string            `json:"path,omitempty"`
	RepoName    string
	SHA         *string            `json:"sha,omitempty" xorm:"sha"`
	HTMLURL     *string            `json:"html_url,omitempty" xorm:"html_url"`
	Repository  *github.Repository `json:"repository,omitempty" xorm:"json"`
	TextMatches []TextMatch        `json:"text_matches,omitempty" xorm:"LONGBLOB"`
	Status      int
	Version     int                `xorm:"version"`
	CreatedTime time.Time          `xorm:"created"`
	UpdatedTime time.Time          `xorm:"updated"`
}

// CodeSearchResult represents the result of a code search.
type CodeSearchResult struct {
	Total             *int         `json:"total_count,omitempty"`
	IncompleteResults *bool        `json:"incomplete_results,omitempty"`
	CodeResults       []CodeResult `json:"items,omitempty" xorm:"json"`
}

func (r *CodeResult) Insert() (int64, error) {
	return Engine.Insert(r)
}

func (r *CodeResult) Exist() (bool, error) {
	var c CodeResult
	return Engine.Table("code_result").Where("name=? and sha=?", r.Name, r.SHA).Get(&c)
}

func ListGithubSearchResult() ([]CodeResult, error) {
	results := make([]CodeResult, 0)
	err := Engine.Where("status=0").Find(&results)
	return results, err
}

func ListGithubSearchResultPage(page int) ([]CodeResult, int, error) {
	results := make([]CodeResult, 0)
	totalPages, err := Engine.Table("code_result").Where("status=0").Count()
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

	err = Engine.Where("status=0").Omit("repository").Desc("id").Limit(vars.PAGE_SIZE,
		(page-1)*vars.PAGE_SIZE).Find(&results)

	return results, pages, err
}

func ListHistoryGithubSearchResultPage(page int) ([]CodeResult, int, error) {
	results := make([]CodeResult, 0)
	totalPages, err := Engine.Table("code_result").Count()
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

	err = Engine.Omit("repository").Desc("id").Limit(vars.PAGE_SIZE,
		(page-1)*vars.PAGE_SIZE).Find(&results)

	return results, pages, err
}


func ListConfirmGithubResultPage(page int) ([]CodeResult, int, error) {
	results := make([]CodeResult, 0)
	totalPages, err := Engine.Table("code_result").Where("status=1").Count()
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

	err = Engine.Where("status=1").Omit("repository").Desc("id").Limit(vars.PAGE_SIZE,
		(page-1)*vars.PAGE_SIZE).Find(&results)

	return results, pages, err
}

func GetReportById(id int64) (bool, *CodeResult, error) {
	report := new(CodeResult)
	has, err := Engine.Id(id).Omit("repository").Get(report)

	return has, report, err
}

func ConfirmReportById(id int64) (err error) {
	report := new(CodeResult)
	has, err := Engine.Id(id).Get(report)
	if err == nil && has {
		report.Status = 1
		_, err = Engine.Id(id).Cols("status").Update(report)
	}
	return err
}

func ResetReportById(id int64) (err error) {
	report := new(CodeResult)
	has, err := Engine.Id(id).Get(report)
	if err == nil && has {
		report.Status = 0
		_, err = Engine.Id(id).Cols("status").Update(report)
	}
	return err
}

func CancelReportById(id int64) (err error) {
	report := new(CodeResult)
	has, err := Engine.Id(id).Omit("repository").Get(report)
	if err == nil && has {
		report.Status = 2
		_, err = Engine.Id(id).Cols("status").Update(report)
		repoName := report.RepoName
		CancelReportByRepo(repoName)
	}
	return err
}

func CancelReportByRepo(repo string) (err error) {

	_, err = Engine.Table("code_result").Exec("update code_result set status=2 where repo_name=?", repo)

	return err
}
