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

package cmd

import (
	"x-patrol/web"
	"x-patrol/util"

	"github.com/urfave/cli"
)

var Web = cli.Command{
	Name:        "web",
	Usage:       "Startup a web Service",
	Description: "Startup a web Service",
	Action:      web.RunWeb,
	Flags: []cli.Flag{
		boolFlag("debug, d", "Debug Mode"),
		stringFlag("host, H", "0.0.0.0", "web listen address"),
		intFlag("port, p", 8000, "web listen port"),
	},
}

var Scan = cli.Command{
	Name:        "scan",
	Usage:       "start to scan github leak info",
	Description: "start to scan github leak info",
	Action:      util.Scan,
	Flags: []cli.Flag{
		stringFlag("mode, m", "github", "scan mode: github, local, all"),
		intFlag("time, t", 900, "scan interval(second)"),
	},
}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func intFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}
