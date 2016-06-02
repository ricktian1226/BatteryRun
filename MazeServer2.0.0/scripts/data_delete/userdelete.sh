#!/bin/sh
if [ $# != 2 ];then
    echo "Usage : $0 <ip:port/dbname> <uid>"
	echo "e.g: $0 192.168.93.129:27017 14124535674528"
	exit 1
fi

urlios="$1/briosdb"
urlandroid="$1/brandroiddb"
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
else
    condition="{uid:\"$uid\"}" 
fi

#account
echo "db.account.remove($condition)" >> $tmpfile
#game
echo "db.game.remove($condition)" >> $tmpfile
#jigsaw
echo "db.jigsaw.remove($condition)" >> $tmpfile
#rune
echo "db.rune.remove($condition)" >> $tmpfile
#consumable
echo "db.consumable.remove($condition)" >> $tmpfile
#usersigninactivity
echo "db.usersigninactivity.remove($condition)" >> $tmpfile
#usermission
echo "db.usermission.remove($condition)" >> $tmpfile
#userdonecollectedmission
echo "db.userdonecollectedmission.remove($condition)" >> $tmpfile
#receipt
echo "db.receipt.remove($condition)" >> $tmpfile
#iaptransaction
echo "db.iaptransaction.remove($condition)" >> $tmpfile
#pushnotification
echo "db.pushnotification.remove($condition)" >> $tmpfile
#pushrecord
echo "db.pushrecord.remove($condition)" >> $tmpfile
#lottosysinfo
echo "db.lottosysinfo.remove($condition)" >> $tmpfile
#lottotransaction
echo "db.lottotransaction.remove($condition)" >> $tmpfile
#shoppingtransaction
echo "db.shoppingtransaction.remove($condition)" >> $tmpfile
#roleinfo
echo "db.roleinfo.remove($condition)" >> $tmpfile
#systemmaillist
echo "db.systemmaillist.remove($condition)" >> $tmpfile
#staminagiveapplylog
echo "db.staminagiveapplylog.remove($condition)" >> $tmpfile
#friendmail
echo "db.friendmail.remove($condition)" >> $tmpfile
#memcache
echo "db.memcache.remove($condition)" >> $tmpfile
$mongoExe $urlios -u $mongouser -p $mongopwd $tmpfile
$mongoExe $urlandroid -u $mongouser -p $mongopwd $tmpfile

