#!/bin/sh
kill -9 $(ps -ef|grep x-patrol-dev|gawk '$0 !~/grep/ {print $2}' |tr -s '\n' ' ')
rm -f ./x-patrol-dev
go build -o x-patrol-dev
./x-patrol-dev