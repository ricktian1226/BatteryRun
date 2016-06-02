#!/bin/sh
if [ $# != 2 ];then
    echo "Usage : $0 <ip:port/dbname> <uid>"
	echo "e.g: $0 192.168.93.129:27017 14124535674528"
	exit 1
fi

urllog="$1/brlogdb"
uid=$2
tmpfile="tmplog.js"
mongoExe=/usr/bin/mongo
mongoUser=superUser
mongoPwd=superUser

if [ -f $tmpfile ];then
    rm $tmpfile
fi

if [ $uid = "all" ];then
    condition="{}"
else
    condition="{uid:\"$uid\"}" 
fi

#accountlog
echo "db.accountlog.remove($condition)" > $tmpfile
#giftlog
echo "db.giftlog.remove($condition)" >> $tmpfile
#gamelog
echo "db.gamelog.remove($condition)" >> $tmpfile
#shoppinglog
echo "db.shoppinglog.remove($condition)" >> $tmpfile
#iaplog
echo "db.iaplog.remove($condition)" >> $tmpfile
#lottolog
echo "db.lottolog.remove($condition)" >> $tmpfile
#checkpointlog
echo "db.checkpointlog.remove($condition)" >> $tmpfile
#operationlog
echo "db.operationlog.remove($condition)" >> $tmpfile

$mongoExe $urllog -u $mongoUser -p $mongoPwd $tmpfile

