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

	"time"
)

type InputInfo struct {
	Id          int64
	Type        string    `xorm:"varchar(255) notnull"`
	Content     string    `xorm:"text notnull"`
	Desc        string    `xorm:"text notnull"`
	Version     int       `xorm:"version"`
	CreatedTime time.Time `xorm:"created"`
	UpdatedTime time.Time `xorm:"updated"`
}

func NewInputInfo(inputType, content, desc string) (info *InputInfo) {
	return &InputInfo{Type: inputType, Content: content, Desc: desc}
}

func (i *InputInfo) Insert() (int64, error) {
	return Engine.Insert(i)
}

func (i *InputInfo) Exist(repoUrl string) (bool, error) {
	info := new(InputInfo)
	return Engine.Table("input_info").Where("content=?", repoUrl).Get(info)
}

func GetInputInfoById(id int64) (*InputInfo, bool, error) {
	input := InputInfo{Id: id}
	has, err := Engine.Table("input_info").ID(id).Get(&input)
	return &input, has, err
}

func EditInputInfoById(id int64, inputType, content, desc string) (error) {
	input := new(InputInfo)
	var err error
	has, err := Engine.ID(id).Get(input)
	if err == nil && has {
		input.Type = inputType
		input.Content = content
		input.Desc = desc
		_, err = Engine.ID(id).Update(input)
	}
	return err
}

func DeleteInputInfoById(id int64) (error) {
	input := new(InputInfo)
	_, err := Engine.Table("input_info").ID(id).Delete(input)
	return err
}

func DeleteAllInputInfo() (error) {
	sqlCMD := "delete from input_info;"
	_, err := Engine.Exec(sqlCMD)
	return err
}

func ListInputInfo() ([]InputInfo, error) {
	inputs := make([]InputInfo, 0)
	err := Engine.Table("input_info").Find(&inputs)
	return inputs, err
}

func ListInputInfoPage(page int) ([]InputInfo, int, error) {
	inputs := make([]InputInfo, 0)

	totalPages, err := Engine.Table("input_info").Count()
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

	err = Engine.Table("input_info").Limit(vars.PAGE_SIZE, (page-1)*vars.PAGE_SIZE).Find(&inputs)

	return inputs, pages, err
}
