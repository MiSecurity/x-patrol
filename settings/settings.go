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

package settings

import (
	"x-patrol/logger"
	"x-patrol/vars"

	"gopkg.in/ini.v1"
)

var (
	Cfg *ini.File
)

func init() {
	var err error
	source := "conf/app.ini"
	Cfg, err = ini.Load(source)

	if err != nil {
		logger.Log.Panicln(err)
	}

	vars.HTTP_HOST = Cfg.Section("").Key("HTTP_HOST").MustString("127.0.0.1")
	vars.HTTP_PORT = Cfg.Section("").Key("HTTP_PORT").MustInt(8000)

	vars.REPO_PATH = Cfg.Section("").Key("REPO_PATH").MustString("repos")
	vars.MAX_INDEXERS = Cfg.Section("").Key("MAX_INDEXERS").MustInt(vars.DefaultMaxConcurrentIndexers)
	vars.MAX_Concurrency_REPOS = Cfg.Section("").Key("MAX_Concurrency_REPOS").MustInt(100)

}
