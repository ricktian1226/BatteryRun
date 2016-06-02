#!/bin/sh
if [ $# != 1 ];then
    echo "Usage : $0 <ip:port>"
	echo "e.g: $0 192.168.93.129:27017"
	exit 1
fi

mongoExe=/usr/bin/mongo
url="$1/brdb02"
tmpfile="tmp.js"
if [ -f $tmpfile ];then
    rm $tmpfile
fi
#account
echo "db.account.remove({})" >> $tmpfile
#tpidmap
echo "db.tpidmap.remove({})" >> $tmpfile
#game
echo "db.tpidmap.remove({})" >> $tmpfile
#jigsaw
echo "db.jigsaw.remove({})" >> $tmpfile
#rune
echo "db.rune.remove({})" >> $tmpfile
#consumable
echo "db.consumable.remove({})" >> $tmpfile
#usersigninactivity
echo "db.usersigninactivity.remove({})" >> $tmpfile
#usermission
echo "db.usermission.remove({})" >> $tmpfile
#userdonecollectedmission
echo "db.userdonecollectedmission.remove({})" >> $tmpfile
#gift
echo "db.gift.remove({})" >> $tmpfile
#receipt
echo "db.receipt.remove({})" >> $tmpfile
#iaptransaction
echo "db.iaptransaction.remove({})" >> $tmpfile
#pushnotification
echo "db.pushnotification.remove({})" >> $tmpfile
#pushrecord
echo "db.pushrecord.remove({})" >> $tmpfile
#announcement
echo "db.announcement.remove({})" >> $tmpfile
#lottosysinfo
echo "db.lottosysinfo.remove({})" >> $tmpfile
#lottotransaction
echo "db.lottotransaction.remove({})" >> $tmpfile
#shoppingtransaction
echo "db.shoppingtransaction.remove({})" >> $tmpfile
#usercheckpoint
echo "db.usercheckpoint.remove({})" >> $tmpfile
#roleinfo
echo "db.roleinfo.remove({})" >> $tmpfile
#systemmaillist
echo "db.systemmaillist.remove({})" >> $tmpfile
#staminagiveapplylog
echo "db.staminagiveapplylog.remove({})" >> $tmpfile
#friendmail
echo "db.friendmail.remove({})" >> $tmpfile
#giftlog
echo "db.giftlog.remove({})" >> $tmpfile
#gamelog
echo "db.gamelog.remove({})" >> $tmpfile
#shoppinglog
echo "db.shoppinglog.remove({})" >> $tmpfile
#accountlog
echo "db.accountlog.remove({})" >> $tmpfile
#iaplog
echo "db.iaplog.remove({})" >> $tmpfile
#lottolog
echo "db.lottolog.remove({})" >> $tmpfile
#checkpointlog
echo "db.checkpointlog.remove({})" >> $tmpfile
#operationlog
echo "db.operationlog.remove({})" >> $tmpfile

$mongoExe $url $tmpfile
