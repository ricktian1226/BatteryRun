【注意】
该目录的脚本都是删除业务动态数据的脚本，慎用！！！

deletedynamicdata.sh 删除服务动态数据。测试用，在重置全网服务数据的时候使用。现网场景不要使用！！！
userdelete.sh        删除某一玩家动态数据。测试用，在删除某一玩家的时候使用。现网场景不要使用！！！
userlogdelete.sh     删除某一玩家日志数据。测试用，在删除某一玩家的时候使用。现网场景不要使用！！！
userkill.sh          删除某一玩家日志数据。测试用，在删除某一玩家的时候使用。现网场景不要使用！！！

目录crontab下的脚本为生产环境定期删除业务动态数据的定时任务：
friendmaildelete.sh     删除过期的玩家好友邮件
gamedelete.sh           删除过期的玩家游戏信息
iaptransactiondelete.sh 删除过期的iap交易信息
lottotransaction.sh     删除过期的抽奖事务信息
shoppingtransactiondelete.sh 删除过期的玩家购买信息
systemmaillist.sh            删除过期的玩家邮件信息
usermissiondelete.sh         删除过期的玩家任务信息
logdelete.sh   删除过期的业务日志信息