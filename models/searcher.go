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
	"x-patrol/util/index"
	"x-patrol/vars"
)

type SearchResult struct {
	Id       int64
	Repo     string
	Filename string
	Match    *index.Match `xorm:"extends"`
	// FilesWithMatch int
	// FilesOpened    int                `json:"-"`
	// Duration       time.Duration      `json:"-"`
	Revision string
	Hash     string
	Rule     Rules
	Status   int `xorm:"int default 0 notnull"`
}

func NewSearchResult(
	match *index.Match,
	repo string,
	filename string,
// FilesWithMatch int,
// FilesOpened int,
// duration time.Duration,
	revision string,
	hash string,
	rule Rules) (*SearchResult) {
	return &SearchResult{Match: match, Repo: repo, Filename: filename, Revision: revision, Hash: hash, Rule: rule}
}

func (s *SearchResult) Insert() (err error) {
	has, err := s.Exist()
	if !has {
		_, err = Engine.Insert(s)
	}

	return err
}

func (s *SearchResult) Exist() (bool, error) {
	result := new(SearchResult)
	return Engine.Table("search_result").Where("repo=? and filename=?", s.Repo, s.Filename).Get(result)

}

func ListSearchResult() ([]SearchResult, error) {
	result := make([]SearchResult, 0)
	err := Engine.Where("status=0").Find(&result)
	return result, err
}

func ListSearchResultPage(page int) ([]SearchResult, int, error) {
	result := make([]SearchResult, 0)
	totalPages, err := Engine.Table("search_result").Where("status=0").Count()
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
	err = Engine.Where("status=0").Limit(vars.PAGE_SIZE, (page-1)*vars.PAGE_SIZE).Find(&result)
	return result, pages, err
}

func GetSearchResultById(id int64) (bool, *SearchResult, error) {
	result := new(SearchResult)
	has, err := Engine.ID(id).Get(result)
	return has, result, err
}

func ConfirmSearchResultById(id int64) (err error) {
	result := new(SearchResult)
	has, err := Engine.ID(id).Get(result)
	if err == nil && has {
		result.Status = 1
		_, err = Engine.ID(id).Update(result)
	}
	return err
}

func CancelSearchResultById(id int64) (err error) {
	result := new(SearchResult)
	has, err := Engine.ID(id).Get(result)
	if err == nil && has {
		result.Status = 2
		_, err = Engine.ID(id).Update(result)
		repoName := result.Repo
		err = CancelSearchResultByRepo(repoName)
	}
	return err
}

func CancelSearchResultByFileName(filename string) (err error) {

	_, err = Engine.Table("search_result").Exec("update search_result set status=2 where filename=?", filename)

	return err
}

func CancelSearchResultByRepo(repo string) (err error) {

	_, err = Engine.Table("search_result").Exec("update search_result set status=2 where repo=?", repo)

	return err
}

func GetRepoUrlById(id int64) (string) {
	result := new(SearchResult)
	Engine.ID(id).Get(result)
	return result.Repo
}
