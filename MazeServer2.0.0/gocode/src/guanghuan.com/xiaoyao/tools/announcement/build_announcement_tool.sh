#!/bin/bash
GOPATH=$HOME/svn/Savage/Program/MazeServer/gocode
#cd $GOPATH
#go build -o announcement_tool -gcflags '-N -l' -race -v ../gocode/src/guanghuan.com/xiaoyao/tools/announcement
go build -o announcement_tool -gcflags '-N -l' -race -v ./
