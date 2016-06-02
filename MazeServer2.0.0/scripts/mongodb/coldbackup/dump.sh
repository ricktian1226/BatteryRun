#!/bin/bash
#数据库备份脚本
#crontab : 0 11,23 * * * /home/mongodb/dump.sh 54.200.192.17:20145
#用法:./$0 <host:port>. host 请求服务器，port 请求端口.
MONGODIR=/home/mongodb/mongodb
#URL=54.200.192.17:20145

$URL=$1
if [ ! -n "$URL" ];then
  echo -e "Usage: $0 <host:port>. \n\t$0 54.200.192.17:20145, for example."
  exit 2
fi

DBNAME=brdb01

DUMPDIR=/home/mongodb/dump
if [ ! -x $DUMPDIR ];then
  mkdir -p $DUMPDIR
fi

LOGDIR=/home/mongodb/log
if [ ! -x $LOGDIR ];then
  mkdir -p $LOGDIR
fi
LOG=$LOGDIR/dump.log


DUMPFILE=`date +%Y%m%d%H%M%S`
`$MONGODIR/bin/mongodump --host $URL -d $DBNAME -o $DUMPDIR/$DUMPFILE >> $LOG`
