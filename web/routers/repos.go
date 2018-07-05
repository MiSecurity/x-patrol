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

package routers

import (
	"x-patrol/models"
	"x-patrol/vars"

	"gopkg.in/macaron.v1"

	"github.com/go-macaron/session"
	"github.com/go-macaron/csrf"

	"strconv"
	"net/url"
	"strings"
)

func ListReposConf(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	if p < 1 {
		p = 1
	}
	pre := p - 1
	if pre <= 0 {
		pre = 1
	}
	next := p + 1

	if sess.Get("admin") != nil {
		repos, pages, _ := models.ListRepoConfigPage(p)
		pList := 0
		if pages-p > 10 {
			pList = p + 10
		} else {
			pList = pages
		}

		pageList := make([]int, 0)
		if pages <= 10 {
			for i := 1; i <= pList; i++ {
				pageList = append(pageList, i)
			}
		} else {
			if p <= 10 {
				for i := 1; i <= pList; i++ {
					pageList = append(pageList, i)
				}
			} else {
				t := p + 5
				if t > pages {
					t = pages
				}
				for i := p - 5; i <= t; i++ {
					pageList = append(pageList, i)
				}
			}
		}

		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["repos"] = repos
		ctx.HTML(200, "repos_conf")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func NewRepoConf(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		ctx.HTML(200, "repo_conf_new")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoNewRepoConf(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		repoName := strings.TrimSpace(ctx.Req.Form.Get("name"))
		repoUrl := strings.TrimSpace(ctx.Req.Form.Get("url"))
		repoType := strings.TrimSpace(ctx.Req.Form.Get("type"))
		username := strings.TrimSpace(ctx.Req.Form.Get("username"))
		password := strings.TrimSpace(ctx.Req.Form.Get("password"))

		_, _, urlPat := models.GetUrlPattern(repoType)
		vcsConf := models.VcsConfig{Username: username, Password: password}

		repoConfig := models.NewRepoConfig(repoName, repoUrl, vars.DefaultPollInterval, repoType, *urlPat, vcsConf,
			false, false)

		has, err := repoConfig.Exist()
		if err == nil && !has {
			repoConfig.Insert()
		}

		ctx.Redirect("/admin/repos/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EditRepoConf(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		assets, _, _ := models.GetInputInfoById(int64(Id))
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["assets"] = assets
		ctx.Data["user"] = sess.Get("admin")
		ctx.HTML(200, "repo_conf_edit")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoEditRepoConf(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		Type := strings.TrimSpace(ctx.Req.Form.Get("type"))
		content := strings.TrimSpace(ctx.Req.Form.Get("content"))
		desc := strings.TrimSpace(ctx.Req.Form.Get("desc"))
		models.EditInputInfoById(int64(Id), Type, content, desc)
		ctx.Redirect("/admin/repos/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EnableRepoConf(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.EnableRepoConfById(int64(Id))
		refer := "/admin/repos/list/"
		if ctx.Req.Header["Referer"] != nil && len(ctx.Req.Header["Referer"]) > 0 {
			u := ctx.Req.Header["Referer"][0]
			urlParsed, err := url.Parse(u)
			if err == nil {
				refer = urlParsed.RequestURI()
			}
		}
		ctx.Redirect(refer)
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DisableRepoConf(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.DisableRepoConfById(int64(Id))
		refer := "/admin/repos/list/"
		if ctx.Req.Header["Referer"] != nil && len(ctx.Req.Header["Referer"]) > 0 {
			u := ctx.Req.Header["Referer"][0]
			urlParsed, err := url.Parse(u)
			if err == nil {
				refer = urlParsed.RequestURI()
			}
		}
		ctx.Redirect(refer)
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteAllRepoConf(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		models.DeleteAllReposConf()
		refer := "/admin/repos/list/"
		if ctx.Req.Header["Referer"] != nil && len(ctx.Req.Header["Referer"]) > 0 {
			u := ctx.Req.Header["Referer"][0]
			urlParsed, err := url.Parse(u)
			if err == nil {
				refer = urlParsed.RequestURI()
			}
		}
		ctx.Redirect(refer)
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DelRepoConfById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.DeleteRepoConfById(int64(Id))
		refer := "/admin/repos/list/"
		ctx.Redirect(refer)
	} else {
		ctx.Redirect("/admin/login/")
	}
}
