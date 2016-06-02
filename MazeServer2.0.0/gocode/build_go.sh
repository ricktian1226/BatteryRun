#!/bin/sh

export GOPATH=/home/xyao/workspace/battery/MazeServer2.0.0/gocode
PROTO_DIR=/home/xyao/workspace/Battery_Interface/MazeNetMsgDef2.0.0/battery_run_net
TARGET_DIR=${GOPATH}/src/guanghuan.com/xiaoyao/superbman_server/battery_run_net
cd $PROTO_DIR
git fetch origin
git reset --hard origin/$1
chmod +x /home/xyao/workspace/Battery_Interface/MazeNetMsgDef2.0.0/protoc
chmod +x /home/xyao/workspace/Battery_Interface/MazeNetMsgDef2.0.0/protoc-gen-go
protoc --go_out=${TARGET_DIR} --proto_path=${PROTO_DIR} ${PROTO_DIR}/*.proto




