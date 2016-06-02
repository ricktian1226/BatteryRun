#!/bin/bash
set GOPATH=$HOME/workspace/MazeServer/gocode
export GOPATH
cd $GOPATH/src/guanghuan.com/xiaoyao
svn update
cd $GOPATH
go build -o battery_server -race -v ./src/guanghuan.com/xiaoyao/superbman_server
