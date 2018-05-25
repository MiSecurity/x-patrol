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

	"os"
	"encoding/json"
	"bufio"
	"io/ioutil"
)

type Rules struct {
	Id          int64
	Part        string
	Type        string
	Pattern     string
	Caption     string
	Description string `xorm:"text"`
	Status      int    `xorm:"int default 0 notnull"`
}

func NewRules(part, ruleType, pat, caption, desc string, status int) (*Rules) {
	return &Rules{Part: part, Type: ruleType, Pattern: pat, Caption: caption, Description: desc, Status: status}
}

func (r *Rules) Insert() (err error) {
	_, err = Engine.Insert(r)
	return err
}

func GetLocalRules() ([]Rules, error) {
	rules := make([]Rules, 0)
	err := Engine.Table("rules").Where("part <> 'github' and status=1").Find(&rules)
	return rules, err
}

func GetAllRules() ([]Rules, error) {
	rules := make([]Rules, 0)
	err := Engine.Table("rules").Find(&rules)
	return rules, err
}

func GetRulesPage(page int) ([]Rules, int, error) {
	rules := make([]Rules, 0)
	totalPages, err := Engine.Table("rules").Count()
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
	err = Engine.Table("rules").Limit(vars.PAGE_SIZE, (page-1)*vars.PAGE_SIZE).Find(&rules)
	return rules, pages, err
}

func GetRuleById(id int64) (*Rules, bool, error) {
	rule := new(Rules)
	has, err := Engine.ID(id).Get(rule)
	return rule, has, err
}

func EditRuleById(id int64, part, ruleType, pat, caption, desc string, status int) (error) {
	rule := new(Rules)
	_, has, err := GetRuleById(id)
	if err == nil && has {
		rule.Part = part
		rule.Type = ruleType
		rule.Pattern = pat
		rule.Caption = caption
		rule.Description = desc
		rule.Status = status
		_, err = Engine.ID(id).Update(rule)
	}
	return err
}

func DeleteRulesById(id int64) (err error) {
	rule := new(Rules)
	_, err = Engine.Id(id).Delete(rule)
	return err
}

func EnableRulesById(id int64) (err error) {
	rules := new(Rules)
	has, err := Engine.Id(id).Get(rules)
	if err == nil && has {
		rules.Status = 1
		_, err = Engine.Id(id).Cols("status").Update(rules)
	}
	return err
}

func DisableRulesById(id int64) (err error) {
	rules := new(Rules)
	has, err := Engine.Id(id).Get(rules)
	if err == nil && has {
		rules.Status = 0
		_, err = Engine.Id(id).Cols("status").Update(rules)
	}
	return err
}

func LoadRuleFromFile(filename string) ([]Rules, error) {
	ruleFile, err := os.Open(filename)
	rules := make([]Rules, 0)
	var content []byte
	if err == nil {
		r := bufio.NewReader(ruleFile)
		content, err = ioutil.ReadAll(r)
		if err == nil {
			err = json.Unmarshal(content, &rules)
		}
	}
	return rules, err
}

func InsertRules(filename string) (error) {
	rules, err := LoadRuleFromFile(filename)
	if err == nil {
		for _, rule := range rules {
			rule.Insert()
		}
	}
	return err
}

func GetGithubKeywords() ([]Rules, error) {
	rules := make([]Rules, 0)
	err := Engine.Table("rules").Where("part='github' and status=1").Find(&rules)
	return rules, err
}
