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
	"github.com/go-macaron/captcha"

	"strings"
	"strconv"
)

func AdminIndex(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		ctx.Redirect("/admin/reports/github/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func Login(ctx *macaron.Context, x csrf.CSRF) {
	ctx.Data["csrf_token"] = x.GetToken()
	ctx.HTML(200, "login")
}

func DoLogin(ctx *macaron.Context, sess session.Store, cpt *captcha.Captcha) {
	if cpt.VerifyReq(ctx.Req) {
		username := ctx.Req.Form.Get("username")
		password := ctx.Req.Form.Get("password")
		has, err := models.Auth(username, password)
		if err == nil && has {
			sess.Set("admin", username)
			ctx.Redirect("/admin/index/")
		} else {
			ctx.Redirect("/admin/login/")
		}
	} else {
		message := "验证码输入错误"
		ctx.Data["message"] = message
		ctx.HTML(200, "error")
	}
}

func ListUsers(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		users, _ := models.ListAdmins()
		ctx.Data["users"] = users
		ctx.HTML(200, "users")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func NewUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		ctx.HTML(200, "users_new")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoNewUser(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		username := strings.TrimSpace(ctx.Req.Form.Get("username"))
		password := strings.TrimSpace(ctx.Req.Form.Get("password"))
		admin := models.NewAdmin(username, password)
		admin.Insert()
		ctx.Redirect("/admin/users/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EditUser(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		admin, _, _ := models.GetAdminById(int64(Id))
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["user"] = admin
		ctx.Data["admin"] = sess.Get("admin")
		ctx.HTML(200, "users_edit")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoEditUser(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		username := strings.TrimSpace(ctx.Req.Form.Get("username"))
		password := strings.TrimSpace(ctx.Req.Form.Get("password"))
		models.EditAdminById(int64(Id), username, password)
		ctx.Redirect("/admin/users/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.DeleteAdminById(int64(Id))
		ctx.Redirect("/admin/users/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
