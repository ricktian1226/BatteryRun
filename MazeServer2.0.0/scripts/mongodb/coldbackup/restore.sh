#!/bin/bash
#数据库备份脚本
#用法:./$0 <host:port> <dumpname>. host 请求服务器，port 请求端口，dumpname 备份的目录名.
#URL=127.0.0.1:27017
URL=$1
DUMPFILE=$2

if [[ ! -n "$URL" ]] || [[ ! -n "$DUMPFILE" ]]; then 
  echo -e "Usage: $0 <host:port> <dumpfile>. \n\t$0 54.200.192.17:20145 20140630160443, for example."
  exit 2
fi

DUMPDIR=/home/mongodb/dump
MONGODIR=/home/mongodb/mongodb
DBNAME=brdb01

LOGDIR=/home/mongodb/log
LOG=$LOGDIR/restore.log
if [ ! -x $LOGDIR ];then
  mkdir -p $LOGDIR
fi

`$MONGODIR/bin/mongorestore --host $URL -d $DBNAME $DUMPDIR/$DUMPFILE/* >> $LOG`
