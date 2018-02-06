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
	"github.com/go-macaron/csrf"

	"strings"
	"strconv"
)

func ListTokens(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		tokens, _ := models.ListTokens()
		ctx.Data["tokens"] = tokens
		ctx.HTML(200, "tokens")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func NewTokens(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		ctx.HTML(200, "tokens_new")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoNewTokens(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		tokens := strings.TrimSpace(ctx.Req.Form.Get("tokens"))
		desc := strings.TrimSpace(ctx.Req.Form.Get("desc"))
		githubToken := models.NewGithubToken(tokens, desc)
		githubToken.Insert()
		ctx.Redirect("/admin/tokens/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EditTokens(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		tokens, _, _ := models.GetTokenById(int64(Id))
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["tokens"] = tokens
		ctx.Data["admin"] = sess.Get("admin")
		ctx.HTML(200, "tokens_edit")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoEditTokens(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		tokens := strings.TrimSpace(ctx.Req.Form.Get("tokens"))
		desc := strings.TrimSpace(ctx.Req.Form.Get("desc"))
		models.EditTokenById(int64(Id), tokens, desc)
		ctx.Redirect("/admin/tokens/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteTokens(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.DeleteTokenById(int64(Id))
		ctx.Redirect("/admin/tokens/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
