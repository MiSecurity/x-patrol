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
	"x-patrol/settings"
	"x-patrol/logger"

	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"

	"path/filepath"
	"fmt"
)

var (
	DATA_TYPE string
	DATA_HOST string
	DATA_PORT int
	DATA_NAME string
	USERNAME  string
	PASSWORD  string
	SSL_MODE  string
	DATA_PATH string

	Engine *xorm.Engine
)

func init() {
	cfg := settings.Cfg
	sec := cfg.Section("database")

	DATA_TYPE = sec.Key("DB_TYPE").MustString("sqlite")
	DATA_HOST = sec.Key("HOST").MustString("127.0.0.1")
	DATA_PORT = sec.Key("PORT").MustInt(3306)
	USERNAME = sec.Key("USER").MustString("username")
	PASSWORD = sec.Key("PASSWD").MustString("password")
	SSL_MODE = sec.Key("SSL_MODE").MustString("disable")
	DATA_PATH = sec.Key("PATH").MustString("db")
	DATA_NAME = sec.Key("NAME").MustString("xsec")

	err := NewDbEngine()
	if err != nil {
		logger.Log.Panicln(err)
	} else {
		Engine.Sync2(new(RepoConfig))
		Engine.Sync2(new(Rules))
		Engine.Sync2(new(SearchResult))
		Engine.Sync2(new(InputInfo))
		Engine.Sync2(new(Admin))
		Engine.Sync2(new(Repo))
		Engine.Sync(new(UrlPattern))
		Engine.Sync2(new(GithubToken))
		Engine.Sync2(new(CodeResult))
		// Engine.Sync2(new(index.Match))

		InitRules()
		InitAdmin()
		InitUrlPattern()

	}
}

// init a database instance
func NewDbEngine() (err error) {
	switch DATA_TYPE {
	//case "sqlite":
	//	cur, _ := filepath.Abs(".")
	//	dataSourceName := fmt.Sprintf("%v/%v/%v.db", cur, DATA_PATH, DATA_NAME)
	//	logger.Log.Infof("sqlite db: %v", dataSourceName)
	//	Engine, err = xorm.NewEngine("sqlite3", dataSourceName)
	//	Engine.Logger().SetLevel(core.LOG_OFF)
	//	err = Engine.Ping()

	case "mysql":
		dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",
			USERNAME, PASSWORD, DATA_HOST, DATA_PORT, DATA_NAME)

		Engine, err = xorm.NewEngine("mysql", dataSourceName)
		Engine.Logger().SetLevel(core.LOG_OFF)
		err = Engine.Ping()
	case "postgres":
		dbSourceName := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", USERNAME, PASSWORD, DATA_HOST,
			DATA_PORT, DATA_NAME, SSL_MODE)
		Engine, err = xorm.NewEngine("postgres", dbSourceName)
		Engine.Logger().SetLevel(core.LOG_OFF)
		err = Engine.Ping()

	default:
		dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8",
			USERNAME, PASSWORD, DATA_HOST, DATA_PORT, DATA_NAME)

		Engine, err = xorm.NewEngine("mysql", dataSourceName)
		Engine.Logger().SetLevel(core.LOG_OFF)
		err = Engine.Ping()
	}

	return err
}

func InitRules() () {
	cur, _ := filepath.Abs(".")
	ruleFile := fmt.Sprintf("%v/conf/gitrob.json", cur)
	rules, err := GetAllRules()
	if err == nil && len(rules) == 0 {
		logger.Log.Infof("Init rules, err: %v", InsertRules(ruleFile))
	}
}

func InitAdmin() {
	admins, err := ListAdmins()
	if err == nil && len(admins) == 0 {
		username := "xsec"
		password := "x@xsec.io"
		admin := NewAdmin(username, password)
		admin.Insert()
	}
}
