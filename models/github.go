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
	"github.com/google/go-github/github"

	"time"
)

type GithubToken struct {
	Id    int64
	Token string
	Desc  string
	// The number of requests per hour the client is currently limited to.
	Limit int `json:"limit"`
	// The number of remaining requests the client can make this hour.
	Remaining int `xorm:"default 5000 notnull" json:"remaining"`
	// The time at which the current rate limit will reset.
	Reset time.Time `json:"reset"`
}

func NewGithubToken(token, desc string) (*GithubToken) {
	return &GithubToken{Token: token, Desc: desc, Limit: 5000, Remaining: 5000}
}

func (g *GithubToken) Insert() (int64, error) {
	return Engine.Insert(g)
}

func (g *GithubToken) Exist() (bool, error) {
	return Engine.Get(g)
}

func ListTokens() ([]GithubToken, error) {
	tokens := make([]GithubToken, 0)
	err := Engine.Find(&tokens)
	return tokens, err
}

func ListValidTokens() ([]GithubToken, error) {
	tokens := make([]GithubToken, 0)
	err := Engine.Table("github_token").Where("remaining>50").Find(&tokens)
	return tokens, err
}

func GetTokenById(id int64) (*GithubToken, bool, error) {
	token := new(GithubToken)
	has, err := Engine.ID(id).Get(token)
	return token, has, err
}

func EditTokenById(id int64, token, desc string) (error) {
	githubToken := new(GithubToken)
	has, err := Engine.ID(id).Get(githubToken)
	if err == nil && has {
		githubToken.Token = token
		githubToken.Desc = desc
		Engine.ID(id).Update(githubToken)
	}
	return err
}

func DeleteTokenById(id int64) (error) {
	token := new(GithubToken)
	_, err := Engine.ID(id).Delete(token)
	return err
}

func UpdateRate(token string, response *github.Response) (error) {
	githubToken := new(GithubToken)
	has, err := Engine.Table("github_token").Where("token=?", token).Get(githubToken)
	if err == nil && has {
		id := githubToken.Id
		githubToken.Remaining = response.Remaining
		githubToken.Reset = response.Reset.Time
		githubToken.Limit = response.Limit
		Engine.ID(id).Update(githubToken)
	}
	return err
}
