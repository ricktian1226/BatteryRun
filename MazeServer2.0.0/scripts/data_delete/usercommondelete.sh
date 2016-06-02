#!/bin/sh
if [ $# != 2 ];then
    echo "Usage : $0 <ip:port/dbname> <uid>"
	echo "e.g: $0 192.168.93.129:27017 14124535674528"
	exit 1
fi

urlcommon="$1/brcommondb"
uid=$2
tmpfile="tmp.js"
mongoExe=/usr/bin/mongo
mongouser=superUser
mongopwd=superUser
if [ -f $tmpfile ];then
    rm $tmpfile
fi

if [ $uid = "all" ];then
    condition="{}"
    condition1="{}"
else
    condition="{uid:\"$uid\"}" 
    condition1="{gid:\"$uid\"}" 
fi

#tpidmap
echo "db.tpidmap.remove($condition1)" > $tmpfile
#useraccomplishment
echo "db.useraccomplishment.remove($condition)" >> $tmpfile
#usercheckpoint
echo "db.usercheckpoint.remove($condition)" >> $tmpfile
$mongoExe $urlcommon -u $mongouser -p $mongopwd $tmpfile

