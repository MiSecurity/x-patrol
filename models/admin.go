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
	"x-patrol/misc"
)

type Admin struct {
	Id       int64
	Username string
	Password string
}

func NewAdmin(username, password string) (*Admin) {
	encryptPass := misc.MakeMd5(password)
	return &Admin{Username: username, Password: encryptPass}
}

func (u *Admin) Insert() (int64, error) {
	return Engine.Insert(u)
}

func ListAdmins() ([]Admin, error) {
	admins := make([]Admin, 0)
	err := Engine.Table("admin").Find(&admins)
	return admins, err
}

func GetAdminById(id int64) (*Admin, bool, error) {
	admin := new(Admin)
	has, err := Engine.ID(id).Get(admin)
	return admin, has, err
}

func EditAdminById(id int64, username, password string) (error) {
	admin := new(Admin)
	has, err := Engine.ID(id).Get(admin)
	if err == nil && has {
		admin.Username = username
		admin.Password = misc.MakeMd5(password)
		Engine.ID(id).Update(admin)
	}
	return err
}

func DeleteAdminById(id int64) (error) {
	admin := new(Admin)
	_, err := Engine.ID(id).Delete(admin)
	return err
}

func Auth(username, password string) (bool, error) {
	admin := new(Admin)
	encryptPass := misc.MakeMd5(password)
	return Engine.Table("admin").Where("username=? and password=?", username, encryptPass).Get(admin)
}
