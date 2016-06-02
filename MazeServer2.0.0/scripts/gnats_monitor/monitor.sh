#!/bin/sh
#* * * * * /home/xiaoyao/server/bin/monitor.sh
#There is some problem with nats reconnection. So when gnats restartes, we should restart all the business service.
WORKDIR="/home/xiaoyao/server"
LOG="$WORKDIR/log/gnatsd/nats_monitor.log"
PID_FILE="$WORKDIR/bin/gnats.pid"
APP_SERVER="battery_app_server"
GATEWAY_SERVER="battery_gateway_server"
GNATS_PID=`ps -ef|grep gnatsd|grep -v grep|grep -v sudo|awk '{print $2}'`
if [ ! -f  $PID_FILE ];then
  echo $GNATS_PID > $PID_FILE
  echo "`date +"%x %X"` create $PID_FILE : $GNATS_PID" >> $LOG
  exit
else
  GNATS_LAST_PID=`cat ${PID_FILE}`
  echo "`date +"%x %X"` pid last($GNATS_LAST_PID),current($GNATS_PID)." >> $LOG
  if [ $GNATS_LAST_PID -ne $GNATS_PID ];then
        echo "`date +"%x %X"` last pid : $GNATS_LAST_PID no equal to pid : $GNATS_PID, kill business server." >> $LOG
        `killall $APP_SERVER`
        `killall $GATEWAY_SERVER`
    echo $GNATS_PID > $PID_FILE
  fi
fi