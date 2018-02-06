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

package githubsearch

import (
	"x-patrol/models"
	"x-patrol/logger"

	"github.com/google/go-github/github"

	"strings"
)

const (
	CONST_REPO  = "repo"
	CONST_REPOS = "repos"
	CONST_ORGS  = "organizations"
	CONST_USER  = "user"
)

func InsertAllRepos() {
	gitClient, _, _ := GetGithubClient()

	assets, err := models.ListInputInfo()
	if err == nil {
		for _, asset := range assets {
			assetType := strings.ToLower(asset.Type)
			name := asset.Content
			switch assetType {
			case CONST_REPO, CONST_REPOS:
				repos := strings.Split(name, ",")
				for _, item := range repos {
					r := models.NewRepo(item, item, &asset)
					has, err := r.Exist()
					if err == nil && !has {
						r.Insert()
					}
				}

			case CONST_ORGS:
				orgs := strings.Split(name, ",")
				var orgsRepos []*github.Repository
				var usersAll []*github.User
				for _, org := range orgs {
					users, resp, err := gitClient.GetOrgsMembers(org)
					usersAll = append(usersAll, users...)
					logger.Log.Println(users, resp, err)
					repos, resp, err := gitClient.GetOrgsRepos(org)
					orgsRepos = append(orgsRepos, repos...)
					models.UpdateRate(gitClient.Token, resp)
				}
				mapRepos := gitClient.GetUsersRepos(usersAll)
				for _, rs := range mapRepos {
					orgsRepos = append(orgsRepos, rs...)
				}

				for _, repo := range orgsRepos {
					r := models.NewRepo(*repo.Name, *repo.HTMLURL, &asset)
					has, err := r.Exist()
					if err == nil && !has {
						r.Insert()
					}
				}

			case CONST_USER:
				var usersRepos []*github.Repository
				users := strings.Split(name, ",")
				mapRepos := gitClient.GetStrUsersRepos(users)
				for _, rs := range mapRepos {
					usersRepos = append(usersRepos, rs...)
				}
				for _, repo := range usersRepos {
					r := models.NewRepo(*repo.Name, *repo.HTMLURL, &asset)
					has, err := r.Exist()
					if err == nil && !has {
						r.Insert()
					}
				}
			}
		}
	}
}
