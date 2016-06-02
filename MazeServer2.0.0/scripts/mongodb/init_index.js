// account 玩家账户信息
db.account.ensureIndex ({uid:1},{unique: true},{background:1});//唯一索引，避免重复插入
// consumable 玩家消耗品背包
db.consumable.ensureIndex({uid:1, random:-1},{background:1});//查询玩家赛前随机道具
db.consumable.ensureIndex({uid:1, id:-1},{background:1});//使用赛前道具
// friendmail 玩家好友邮件
db.friendmail.ensureIndex({uid:1, createtime:-1},{background:1});//1、删除好友邮件（超过固定时间期限）2、获取玩家好友邮件（最近30封）
db.friendmail.ensureIndex({uid:1, friendid:1},{background:1});//1、判断好友邮件是否存在 2、删除好友邮件
db.friendmail.ensureIndex({uid:1, mailtype:1},{background:1});//1、获取所有赠送体力邮件数目（增加体力）2、删除所有体力赠送邮件（一键确认）
// friendmailcount 玩家好友邮件
db.friendmailcount.ensureIndex({uid:1},{background:1});//获取指定玩家好友邮件总数
// game 玩家游戏数据信息
db.game.ensureIndex({id:1}, {unique:true},{background:1});//根据gameid和uid查询玩家游戏数据信息，用于校验游戏的合法性
// iaptransaction iap交易信息
db.iaptransaction.ensureIndex({transactionid:-1},{background:1}); 
// jigsaw 玩家拼图信息
db.jigsaw.ensureIndex({uid:1,id:-1},{background:1});//1、查询玩家拼图数 2、查询玩家对应的拼图数据是否存在
// lottosysinfo 玩家系统抽奖信息
db.lottosysinfo.ensureIndex({uid:1},{unique:true},{background:1});//查询玩家系统抽奖信息
// lottotransaction 玩家抽奖事务信息
db.lottotransaction.ensureIndex({uid:1, lottoid:-1, parentlottoid:-1},{background:1});//查询玩家系统抽奖信息
// receipt 玩家iap账单信息（无查询、修改需求，无需索引）
// roleinfo 玩家角色信息
db.roleinfo.ensureIndex({uid:1},{background:1});//1、查询玩家角色信息 2、更新玩家角色信息
// rune 玩家符文信息
db.rune.ensureIndex({uid:1},{background:1});//查询玩家符文信息
// shoppingtransaction 玩家购买事务信息
db.shoppingtransaction.ensureIndex({uid:1, gameid:-1, goodsid:-1},{background:1});//查询玩家商品交易次数（针对存在每人购买上限和每局购买上限的商品）
// staminagiveapplylog 玩家体力赠送/请求时间戳信息
db.staminagiveapplylog.ensureIndex({uid:1, frienduid:1},{background:1});//查询玩家体力赠送/请求时间戳信息
// systemmail 玩家系统邮件信息
db.systemmaillist.ensureIndex({uid:1},{background:1});//查询玩家系统邮件信息
// usermission 玩家任务信息
db.usermission.ensureIndex({uid:1, type:1, state:1},{background:1});//查询玩家指定类型任务信息
db.usermission.ensureIndex({uid:1, mid:1, state:1},{background:1});//查询玩家指定类型任务信息
db.userdonecollectedmission.ensureIndex({uid:1, type:1},{background:1});//查询玩家指定类型doneCollected任务信息
// usersignactivity 玩家签到活动信息
db.usersignactivity.ensureIndex({uid:1, id:1},{background:1});//查询玩家指定类型任务信息
// usershareinfo 玩家分享信息
db.usershareinfo.ensureIndex({uid:1},{background:1});//查询玩家分享信息