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
	"gopkg.in/macaron.v1"

	"github.com/go-macaron/session"

	"strconv"
	"net/url"
)

func ListGithubSearchResult(ctx *macaron.Context, sess session.Store) {
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
		reports, pages, _ := models.ListGithubSearchResultPage(p)
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

		ctx.Data["reports"] = reports
		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.HTML(200, "report_github")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func ListHistoryGithubSearchResult(ctx *macaron.Context, sess session.Store) {
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
		reports, pages, _ := models.ListHistoryGithubSearchResultPage(p)
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

		ctx.Data["reports"] = reports
		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.HTML(200, "report_github_history")
	} else {
		ctx.Redirect("/admin/login/")
	}
}


func ListConfirmGithubSearchResult(ctx *macaron.Context, sess session.Store) {
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
		reports, pages, _ := models.ListConfirmGithubResultPage(p)
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

		ctx.Data["reports"] = reports
		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.HTML(200, "report_github_confirm")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func ConfirmReportById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		refer := "/admin/reports/github/"
		if ctx.Req.Header["Referer"] != nil && len(ctx.Req.Header["Referer"]) > 0 {
			u := ctx.Req.Header["Referer"][0]
			urlParsed, err := url.Parse(u)
			if err == nil {
				refer = urlParsed.RequestURI()
			}
		}
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.ConfirmReportById(int64(Id))
		ctx.Redirect(refer)
	} else {
		ctx.Redirect("/admin/login/")
	}
}


func ResetReportById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		refer := "/admin/reports/github/"
		if ctx.Req.Header["Referer"] != nil && len(ctx.Req.Header["Referer"]) > 0 {
			u := ctx.Req.Header["Referer"][0]
			urlParsed, err := url.Parse(u)
			if err == nil {
				refer = urlParsed.RequestURI()
			}
		}
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.ResetReportById(int64(Id))
		ctx.Redirect(refer)
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func CancelReportById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		refer := "/admin/reports/github/"
		if ctx.Req.Header["Referer"] != nil && len(ctx.Req.Header["Referer"]) > 0 {
			u := ctx.Req.Header["Referer"][0]
			urlParsed, err := url.Parse(u)
			if err == nil {
				refer = urlParsed.RequestURI()
			}
		}

		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.CancelReportById(int64(Id))
		ctx.Redirect(refer)
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DisableRepoById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)

		refer := "/admin/reports/github/"
		if ctx.Req.Header["Referer"] != nil && len(ctx.Req.Header["Referer"]) > 0 {
			u := ctx.Req.Header["Referer"][0]
			urlParsed, err := url.Parse(u)
			if err == nil {
				refer = urlParsed.RequestURI()
			}
		}

		has, result, err := models.GetReportById(int64(Id))
		if err == nil && has {
			models.DisableRepoByUrl(result.Repository.GetHTMLURL())
		}
		models.CancelReportById(int64(Id))
		ctx.Redirect(refer)
	} else {
		ctx.Redirect("/admin/login/")
	}
}

/*
For local code search
*/

func ListLocalSearchResult(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		reports, _ := models.ListSearchResult()
		ctx.Data["reports"] = reports
		ctx.HTML(200, "report_search")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func ListLocalSearchResultPage(ctx *macaron.Context, sess session.Store) {
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
		reports, pages, _ := models.ListSearchResultPage(p)

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
		ctx.Data["reports"] = reports
		ctx.HTML(200, "report_search")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func ConfirmSearchResultById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.ConfirmSearchResultById(int64(Id))
		refer := "/admin/reports/search/"
		if ctx.Req.Header["Referer"] != nil && len(ctx.Req.Header["Referer"]) > 0 {
			u := ctx.Req.Header["Referer"][0]
			urlParsed, err := url.Parse(u)
			if err == nil {
				refer = urlParsed.RequestURI()
			}
		}
		ctx.Redirect(refer)
		//ctx.HTML(200, "back")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func CancelSearchResultById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.CancelSearchResultById(int64(Id))

		refer := "/admin/reports/search/"
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

func DisableSearchRepoById(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		has, result, err := models.GetSearchResultById(int64(Id))
		if err == nil && has {
			repoUrl := result.Repo
			models.DisableRepoByUrl(repoUrl)
		}
		models.CancelSearchResultById(int64(Id))
		refer := "/admin/reports/search/"
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
