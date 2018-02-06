#!/bin/bash

go build main.go
mv main x-patrol_darwin_amd64

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
mv main x-patrol_linux_amd64