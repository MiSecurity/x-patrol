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

package util

import (
	"x-patrol/tasks"
	"x-patrol/util/githubsearch"
	"x-patrol/logger"

	"github.com/urfave/cli"

	"strings"
	"time"
)

func Scan(ctx *cli.Context) () {
	var ScanMode = "github"
	var Interval time.Duration = 900

	if ctx.IsSet("mode") {
		ScanMode = strings.ToLower(ctx.String("mode"))
	}

	if ctx.IsSet("time") {
		Interval = time.Duration(ctx.Int("time"))
	}

	switch ScanMode {
	case "github":
		logger.Log.Println("scan github code")
		githubsearch.ScheduleTasks(Interval)
	case "local":
		logger.Log.Println("scan local repos")
		tasks.ScheduleTasks(Interval)
	case "all":
		logger.Log.Println("scan github code and local repos")
		go githubsearch.ScheduleTasks(Interval)
		tasks.ScheduleTasks(Interval)
	default:
		logger.Log.Println("scan github code ")
		go githubsearch.ScheduleTasks(Interval)
	}
}
