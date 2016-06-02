#!/bin/sh
# crontab:
# 0 2 * * * /server/bin/crontab/iaptransactiondelete.sh m.br2:20144
if [ $# != 1 ];then
    echo "Usage : $0 <ip:port/dbname>"
	echo "e.g: $0 192.168.93.129:27017"
	exit 1
fi
mongoExe="/usr/bin/mongo"
urlios="$1/briosdb"
urlandroid="$1/brandroiddb"
ltimestamp=`date -d "30 days ago" +%s`
module="iaptransaction"
property="timestamp"
logfile="/logs/crontab/crontab$module.log"
tmpfile="tmp$module.js"
if [ -f $tmpfile ];then
    rm $tmpfile
fi

echo "db.$module.remove({$property:{\$lt:$ltimestamp}})" >> $tmpfile
echo "`date +%Y%m%d%H%M%S` $module delete begin" >> $logfile
$mongoExe $urlios $tmpfile
$mongoExe $urlandroid $tmpfile
echo "`date +%Y%m%d%H%M%S` $module delete end" >> $logfile

