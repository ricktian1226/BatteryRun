#!/bin/bash

if [ $# != 0 ];then
  echo "version is $1"
  dir="$1/"
  mkdir -p $dir
fi 

export GOPATH=/home/xyao/workspace/battery/MazeServer2.0.0/gocode
cd $GOPATH
chmod +x build_go.sh
./build_go.sh $1
git fetch origin
git reset --hard origin/$1
echo "=========== building battery gateway server ============"
go build -o "$dir"battery_gateway_server -v ./src/guanghuan.com/xiaoyao/battery_gateway_server
echo "=========== building battery app server ============"
go build -o "$dir"battery_app_server -v ./src/guanghuan.com/xiaoyao/battery_app_server
echo "=========== building battery apns server ============"
go build -o "$dir"battery_apns_server -v ./src/guanghuan.com/xiaoyao/battery_apns_server
echo "=========== building battery file server ============"
go build -o "$dir"battery_file_server -v ./src/guanghuan.com/xiaoyao/battery_file_server
echo "=========== building battery transaction server ============"
go build -o "$dir"battery_transaction_server -v ./src/guanghuan.com/xiaoyao/battery_transaction_server
echo "=========== building battery mail server ============"
go build -o "$dir"battery_mail_server -v ./src/guanghuan.com/xiaoyao/battery_mail_server
echo "=========== building battery maintenance server ============"
go build -o "$dir"battery_maintenance_server -v ./src/guanghuan.com/xiaoyao/battery_maintenance_server
echo "=========== building battery statistic server ============"
go build -o "$dir"battery_statistic_server -v ./src/guanghuan.com/xiaoyao/battery_statistic_server
cd $dir
rm -f server.tar.gz 
tar zcvf server.tar.gz battery_*_server
