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
)

type Repo struct {
	Id     int64
	Name   string
	Url    string
	Src    *InputInfo `xorm:"json"`
	Status int        `xorm:"int notnull default(1)"`
}

func NewRepo(name, repoUrl string, src *InputInfo) (repo *Repo) {
	return &Repo{Name: name, Url: repoUrl, Src: src, Status: 1}
}

func (r *Repo) Insert() (int64, error) {
	return Engine.Insert(r)
}

func (r *Repo) Exist() (bool, error) {
	repo := new(Repo)
	repo.Name = r.Name
	return Engine.Table("repo").Get(repo)
}

func ListRepos() ([]Repo, error) {
	repos := make([]Repo, 0)
	err := Engine.Find(&repos)
	return repos, err
}

func ListReposPage(page int) ([]Repo, int, error) {
	repos := make([]Repo, 0)
	totalPages, err := Engine.Table("repo").Count()
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

	err = Engine.Limit(vars.PAGE_SIZE, (page-1)*vars.PAGE_SIZE).Find(&repos)
	return repos, pages, err
}

func ListEnableRepos() ([]Repo, error) {
	repos := make([]Repo, 0)
	err := Engine.Table("repo").Where("status=?", 1).Find(&repos)
	return repos, err
}

func EnableRepoById(id int64) (error) {
	repo := new(Repo)
	has, err := Engine.ID(id).Get(repo)
	if err == nil && has {
		repo.Status = 1
		_, err = Engine.ID(id).Cols("status").Update(repo)
	}
	return err
}

func DisableRepoById(id int64) (error) {
	repo := new(Repo)
	has, err := Engine.ID(id).Get(repo)
	if err == nil && has {
		repo.Status = 0
		_, err = Engine.ID(id).Cols("status").Update(repo)
	}
	return err
}

func DeleteAllRepos() (error) {
	sqlCmd := "delete from repos"
	_, err := Engine.Exec(sqlCmd)
	return err
}

func DisableRepoByUrl(repoUrl string) (error) {
	repo := new(Repo)
	has, err := Engine.Table("repo").Where("url=?", repoUrl).Get(repo)
	if err == nil && has {
		repo.Status = 0
		_, err = Engine.ID(repo.Id).Cols("status").Update(repo)
	}
	return err
}
