#!/bin/bash

# 项目放到$GOPATH/src/目录下，然后安装以下依赖
go get github.com/etsy/hound/codesearch/sparse
go get github.com/go-macaron/cache
go get github.com/go-macaron/captcha
go get github.com/go-macaron/csrf
go get github.com/go-macaron/session
go get github.com/go-sql-driver/mysql
go get github.com/go-xorm/core
go get github.com/go-xorm/xorm
go get github.com/google/go-github/github
go get github.com/lib/pq
go get github.com/sirupsen/logrus
go get github.com/urfave/cli
go get github.com/x-cray/logrus-prefixed-formatter
go get golang.org/x/oauth2
go get gopkg.in/ini.v1
go get gopkg.in/macaron.v1

# MAC平台下的编译命令参考
go build main.go
mv main x-patrol_darwin_amd64

# 在mac下交叉编译LINUX版本参考命令
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
mv main x-patrol_linux_amd64
