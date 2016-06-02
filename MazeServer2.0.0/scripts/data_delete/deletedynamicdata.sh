#!/bin/sh
if [ $# != 2 ];then
  echo "Usage : $0 <ip:port> [data|log]"
  echo "e.g: $0 192.168.93.129:27017 data"
  exit 1
fi

mongoExe=/home/software/mongodb/bin/mongo

function delete_file()
{  
  if [ -f $1 ];then
    rm $1
  fi
}


function delete_common_data_func()
{
  tmpfile="tmp$2.js"
  url="$1/$2"
  delete_file $tmpfile
  #echo "use $2" > $tmpfile
  #tpidmap
  echo "db.tpidmap.remove({})" >> $tmpfile
  #usercheckpoint
  echo "db.usercheckpoint.remove({})" >> $tmpfile
  #useraccomplishment
  echo "db.useraccomplishment.remove({})" >> $tmpfile
  
  $mongoExe $url $tmpfile
}

function delete_platform_data_func
{
  tmpfile="tmp$2.js"
  url="$1/$2"
  delete_file $tmpfile
  #echo "use $2" > $tmpfile
  #account
  echo "db.account.remove({})" >> $tmpfile
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
  #echo "db.announcement.remove({})" >> $tmpfile
  #lottosysinfo
  echo "db.lottosysinfo.remove({})" >> $tmpfile
  #lottotransaction
  echo "db.lottotransaction.remove({})" >> $tmpfile
  #shoppingtransaction
  echo "db.shoppingtransaction.remove({})" >> $tmpfile
  #roleinfo
  echo "db.roleinfo.remove({})" >> $tmpfile
  #systemmaillist
  echo "db.systemmaillist.remove({})" >> $tmpfile
  #staminagiveapplylog
  echo "db.staminagiveapplylog.remove({})" >> $tmpfile
  #friendmail
  echo "db.friendmail.remove({})" >> $tmpfile
  #memcache
  echo "db.memcache.remove({})" >> $tmpfile
  
  $mongoExe $url $tmpfile
}

function delete_log_func
{
  tmpfile="tmp$2.js"
  url="$1/$2"
  delete_file $tmpfile
  #echo "use $2" > $tmpfile

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
}

if [ $2 == "data" ];then
  delete_common_data_func "$1" "brcommondb"
  delete_platform_data_func "$1" "briosdb"
  delete_platform_data_func "$1" "brandroiddb"
elif [ $2 == "log" ];then
  delete_log_func "$1" "brlogdb"
fi


