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

package vars

import "time"

const (
	DefaultPollInterval          = 900
	DefaultMaxConcurrentIndexers = 2
	DefaultPollEnabled           = false
	DefaultVcs                   = "git"
	DefaultBaseUrl               = "{url}/blob/master/{path}{anchor}"
	DefaultAnchor                = "#L{line}"
)

var (
	REPO_PATH    string
	MAX_INDEXERS int

	HTTP_HOST string
	HTTP_PORT int

	MAX_Concurrency_REPOS int

	DEBUG_MODE bool

	PAGE_SIZE = 10

	TIME_OUT time.Duration = 60 * 4

	Exts map[string]bool
)
