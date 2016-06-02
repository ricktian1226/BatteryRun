#!/bin/sh
if [ $# != 2 ];then
    echo "Usage : $0 <ip:port> <uid>"
	echo "e.g: $0 192.168.93.129:27017 14124535674528"
	exit 1
fi

mongoExe=/usr/bin/mongo
url="$1/brdb02"
uid=$2
tmpfile="tmp.js"
if [ -f $tmpfile ];then
    rm $tmpfile
fi
#account
echo "db.account.remove({uid:\"$uid\"})" >> $tmpfile
#tpidmap
echo "db.tpidmap.remove({gid:\"$uid\"})" >> $tmpfile
#game
echo "db.tpidmap.remove({uid:\"$uid\"})" >> $tmpfile
#jigsaw
echo "db.jigsaw.remove({uid:\"$uid\"})" >> $tmpfile
#rune
echo "db.rune.remove({uid:\"$uid\"})" >> $tmpfile
#consumable
echo "db.consumable.remove({uid:\"$uid\"})" >> $tmpfile
#usersigninactivity
echo "db.usersigninactivity.remove({uid:\"$uid\"})" >> $tmpfile
#usermission
echo "db.usermission.remove({uid:\"$uid\"})" >> $tmpfile
#userdonecollectedmission
echo "db.userdonecollectedmission.remove({uid:\"$uid\"})" >> $tmpfile
#gift
echo "db.gift.remove({uid:\"$uid\"})" >> $tmpfile
#receipt
echo "db.receipt.remove({uid:\"$uid\"})" >> $tmpfile
#iaptransaction
echo "db.iaptransaction.remove({uid:\"$uid\"})" >> $tmpfile
#pushnotification
echo "db.pushnotification.remove({uid:\"$uid\"})" >> $tmpfile
#pushrecord
echo "db.pushrecord.remove({uid:\"$uid\"})" >> $tmpfile
#announcement
echo "db.announcement.remove({uid:\"$uid\"})" >> $tmpfile
#lottosysinfo
echo "db.lottosysinfo.remove({uid:\"$uid\"})" >> $tmpfile
#lottotransaction
echo "db.lottotransaction.remove({uid:\"$uid\"})" >> $tmpfile
#shoppingtransaction
echo "db.shoppingtransaction.remove({uid:\"$uid\"})" >> $tmpfile
#usercheckpoint
echo "db.usercheckpoint.remove({uid:\"$uid\"})" >> $tmpfile
#roleinfo
echo "db.roleinfo.remove({uid:\"$uid\"})" >> $tmpfile
#systemmaillist
echo "db.systemmaillist.remove({uid:\"$uid\"})" >> $tmpfile
#staminagiveapplylog
echo "db.staminagiveapplylog.remove({uid:\"$uid\"})" >> $tmpfile
#friendmail
echo "db.friendmail.remove({uid:\"$uid\"})" >> $tmpfile
#giftlog
echo "db.giftlog.remove({uid:\"$uid\"})" >> $tmpfile
#gamelog
echo "db.gamelog.remove({uid:\"$uid\"})" >> $tmpfile
#shoppinglog
echo "db.shoppinglog.remove({uid:\"$uid\"})" >> $tmpfile
#accountlog
echo "db.accountlog.remove({uid:\"$uid\"})" >> $tmpfile
#iaplog
echo "db.iaplog.remove({uid:\"$uid\"})" >> $tmpfile
#lottolog
echo "db.lottolog.remove({uid:\"$uid\"})" >> $tmpfile
#checkpointlog
echo "db.checkpointlog.remove({uid:\"$uid\"})" >> $tmpfile
#operationlog
echo "db.operationlog.remove({uid:\"$uid\"})" >> $tmpfile

$mongoExe $url $tmpfile