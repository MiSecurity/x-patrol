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

package web

import (
	"x-patrol/web/routers"
	"x-patrol/logger"
	"x-patrol/vars"

	"github.com/go-macaron/cache"
	"github.com/go-macaron/captcha"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"

	"github.com/urfave/cli"

	"net/http"
	"fmt"
	"runtime"
	"html/template"
	"strings"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func RunWeb(ctx *cli.Context) {

	if ctx.IsSet("debug") {
		vars.DEBUG_MODE = ctx.Bool("debug")
	}
	if ctx.IsSet("host") {
		vars.HTTP_HOST = ctx.String("host")
	}
	if ctx.IsSet("port") {
		vars.HTTP_PORT = ctx.Int("port")
	}

	m := macaron.Classic()

	m.Use(macaron.Renderer(
		macaron.RenderOptions{
			Directory:  "templates",
			Extensions: []string{".tmpl", ".html"},
			Funcs: []template.FuncMap{map[string]interface{}{
				"Replace": func(str *string) string {
					t := strings.Replace(*str, "\n", "<br>", -1)
					return t
				},
				"Split": func(str *string) []string {
					return strings.Split(*str, ",")
				},
				"unescaped": func(x string) interface{} { return template.HTML(x) },
			}},
			Delims:          macaron.Delims{"{{", "}}"},
			Charset:         "UTF-8",
			IndentJSON:      true,
			IndentXML:       true,
			PrefixJSON:      []byte("macaron"),
			PrefixXML:       []byte("macaron"),
			HTMLContentType: "text/html",
		}))

	m.Use(session.Sessioner())
	m.Use(csrf.Csrfer())
	m.Use(cache.Cacher())

	m.Use(captcha.Captchaer(captcha.Options{
		URLPrefix:        "/captcha/",
		FieldIdName:      "captcha_id",
		FieldCaptchaName: "captcha",
		ChallengeNums:    4,
		Width:            150,
		Height:           40,
		Expiration:       600,
		CachePrefix:      "captcha_",
	}))

	m.Get("/", routers.Index)

	m.Group("/admin", func() {
		m.Get("", routers.AdminIndex)
		m.Get("/index/", routers.AdminIndex)
		m.Get("/login/", routers.Login)
		m.Post("/login/", routers.DoLogin)

		m.Group("/users/", func() {
			m.Get("", routers.ListUsers)
			m.Get("/list/", routers.ListUsers)
			m.Get("/new/", routers.NewUser)
			m.Post("/new/", routers.DoNewUser)
			m.Get("/edit/:id", routers.EditUser)
			m.Post("/edit/:id", routers.DoEditUser)
			m.Get("/del/:id", routers.DeleteUser)
		})

		/*m.Group("/assets/", func() {
			m.Get("", routers.ListAssets)
			m.Get("/list/", routers.ListAssets)
			m.Get("/list/:page", routers.ListAssets)
			m.Get("/new/", routers.NewAssets)
			m.Post("/new/", routers.DoNewAssets)
			m.Get("/edit/:id", routers.EditAssets)
			m.Post("/edit/:id", routers.DoEditAssets)
			m.Get("/del/:id", routers.DeleteAssets)
			m.Get("/del_all/", routers.DeleteAllAssets)
		})*/

		m.Group("/tokens/", func() {
			m.Get("", routers.ListTokens)
			m.Get("/list/", routers.ListTokens)
			m.Get("/new/", routers.NewTokens)
			m.Post("/new/", routers.DoNewTokens)
			m.Get("/edit/:id", routers.EditTokens)
			m.Post("/edit/:id", routers.DoEditTokens)
			m.Get("/del/:id", routers.DeleteTokens)
		})

		m.Group("/rules/", func() {
			m.Get("", routers.ListRules)
			m.Get("/list/", routers.ListRules)
			m.Get("/list/:page", routers.ListRules)
			m.Get("/new/", routers.NewRules)
			m.Post("/new/", routers.DoNewRules)
			m.Get("/edit/:id", routers.EditRules)
			m.Post("/edit/:id", routers.DoEditRules)
			m.Get("/del/:id", routers.DeleteRules)
			m.Get("/enable/:id", routers.EnableRules)
			m.Get("/disable/:id", routers.DisableRules)
		})

		m.Group("/repos/", func() {
			m.Get("", routers.ListReposConf)
			m.Get("/list/", routers.ListReposConf)
			m.Get("/list/:page", routers.ListReposConf)
			m.Get("/new/", routers.NewRepoConf)
			m.Post("/new/", routers.DoNewRepoConf)
			m.Get("/edit/:id", routers.EditRepoConf)
			m.Post("/edit/:id", routers.DoEditRepoConf)
			m.Get("/enable/:id", routers.EnableRepoConf)
			m.Get("/disable/:id", routers.DisableRepoConf)
			m.Get("/del/:id", routers.DelRepoConfById)
		})

		m.Group("/reports/", func() {
			m.Get("/github/", routers.ListGithubSearchResult)
			m.Get("/github/:page", routers.ListGithubSearchResult)
			m.Get("/github/confirm/:id", routers.ConfirmReportById)
			m.Get("/github/reset/:id", routers.ResetReportById)
			m.Get("/github/cancel/:id", routers.CancelReportById)
			m.Get("/github/disable_repo/:id", routers.DisableRepoById)

			m.Get("/history/", routers.ListHistoryGithubSearchResult)
			m.Get("/history/:page", routers.ListHistoryGithubSearchResult)

			m.Get("/confirm/", routers.ListConfirmGithubSearchResult)
			m.Get("/confirm/:page", routers.ListConfirmGithubSearchResult)

			/* For local repos search */
			m.Get("/search/", routers.ListLocalSearchResultPage)
			m.Get("/search/:page", routers.ListLocalSearchResultPage)
			m.Get("/search/confirm/:id", routers.ConfirmSearchResultById)
			m.Get("/search/cancel/:id", routers.CancelSearchResultById)
			m.Get("/search/disable_repo/:id", routers.DisableSearchRepoById)
		})
	})

	logger.Log.Printf("Server is running on %s", fmt.Sprintf("%v:%v", vars.HTTP_HOST, vars.HTTP_PORT))
	logger.Log.Println(http.ListenAndServe(fmt.Sprintf("%v:%v", vars.HTTP_HOST, vars.HTTP_PORT), m))
}
