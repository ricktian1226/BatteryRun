#!/bin/sh
# crontab:
# 15 3 * * * /server/bin/logdelete.sh a.br2:20145
if [ $# != 1 ];then
    echo "Usage : $0 <ip:port/dbname>"
	echo "e.g: $0 192.168.93.129:27017"
	exit 1
fi
mongoExe="/usr/bin/mongo"
urllog="$1/brlogdb"
logfile="/logs/crontab/crontablog.log"
ltimestamp=`date -d "5 days ago" +%s`
tmpfile="tmplog.js"
if [ -f $tmpfile ];then
    rm $tmpfile
fi
#gamelog
echo "db.gamelog.remove({timestamp:{\$lt:$ltimestamp}})" > $tmpfile
echo "`date +%Y%m%d%H%M%S` gamelog delete begin" >> $logfile
$mongoExe $urllog $tmpfile
echo "`date +%Y%m%d%H%M%S` gamelog delete end" >> $logfile

#shoppinglog
echo "db.shoppinglog.remove({timestamp:{\$lt:$ltimestamp}})" > $tmpfile
echo "`date +%Y%m%d%H%M%S` shoppinglog delete begin" >> $logfile
$mongoExe $urllog $tmpfile
echo "`date +%Y%m%d%H%M%S` shoppinglog delete end" >> $logfile

#accountlog
echo "db.accountlog.remove({timestamp:{\$lt:$ltimestamp}})" > $tmpfile
echo "`date +%Y%m%d%H%M%S` accountlog delete begin" >> $logfile
$mongoExe $urllog $tmpfile
echo "`date +%Y%m%d%H%M%S` accountlog delete end" >> $logfile

#iaplog
echo "db.iaplog.remove({timestamp:{\$lt:$ltimestamp}})" > $tmpfile
echo "`date +%Y%m%d%H%M%S` iaplog delete begin" >> $logfile
$mongoExe $urllog $tmpfile
echo "`date +%Y%m%d%H%M%S` iaplog delete end" >> $logfile

#lottolog
echo "db.lottolog.remove({timestamp:{\$lt:$ltimestamp}})" > $tmpfile
echo "`date +%Y%m%d%H%M%S` lottolog delete begin" >> $logfile
$mongoExe $urllog $tmpfile
echo "`date +%Y%m%d%H%M%S` lottolog delete end" >> $logfile

#pushlog
echo "db.pushlog.remove({timestamp:{\$lt:$ltimestamp}})" > $tmpfile
echo "`date +%Y%m%d%H%M%S` pushlog delete begin" >> $logfile
$mongoExe $urllog $tmpfile
echo "`date +%Y%m%d%H%M%S` pushlog delete end" >> $logfile

#operationlog
echo "db.operationlog.remove({timestamp:{\$lt:$ltimestamp}})" > $tmpfile
echo "`date +%Y%m%d%H%M%S` operationlog delete begin" >> $logfile
$mongoExe $urllog $tmpfile
echo "`date +%Y%m%d%H%M%S` operationlog delete end" >> $logfile

