package batterydb

import (
//	"log"
)

// 数据库中的表
//const (
//	DB_TABLE_ACCOUNT  = "account" // 账户
//	DB_TABLE_TPID_MAP = "tpidmap" // 账户对应表(服务器id与第三方id的映射)
//	//DB_TABLE_FRIENDSHIP = "friendship" // 好友关系
//	DB_TABLE_GAMEDATA = "game" // 每次游戏的成绩
//	//DB_TABLE_STAMINA    = "stamina"    // 玩家当前体力
//	//DB_TABLE_CHIP          = "chip"         // 玩家当前碎片
//	DB_TABLE_JIGSAW        = "jigsaw"       // 玩家当前拼图
//	DB_TABLE_JIGSAW_CONFIG = "jigsawconfig" // 拼图配置表
//	DB_TABLE_RUNE          = "rune"         // 玩家当前符文
//	//DB_TABLE_DAILY_ADWARD  = "dailyaward"   // 每日登录奖励
//	//DB_TABLE_MISSION_LOG   = "missionlog"   // 任务日志
//	DB_TABLE_MISSION     = "missionconfig" // 任务配置信息
//	DB_TABLE_TRANSACTION = "transaction"   // 交易事务
//	//DB_TABLE_AUDIT       = "audit"       // 审计 (各种操作)
//	//DB_TABLE_GIFT      = "gift"     // 体力请求
//	//DB_TABLE_GIFT_DONE = "giftdone" // 已经完成的体力请求
//	DB_TABLE_GOODS   = "gooddsconfig" // 商品定义
//	DB_TABLE_RECEIPT = "receipt"      // 购买凭据/日志
//	//DB_TABLE_DIAMOND = "diamond" // 钻石详细数据
//	//DB_TABLE_IAPORDER     = "iaporder"         // 订单记录
//	DB_TABLE_IAPTRANSACTION = "iaptransaction" // 交易记录
//	//DB_TABLE_DIAMOND_RCD              = "diamondrecord"          // 宝石消耗和产出记录
//	//DB_TABLE_IAPGOOD                  = "iapgood"                // iap购买商品项目
//	DB_TABLE_PUSH_Natice     = "pushnotification"    // 推送命令记录表
//	DB_TABLE_PUSH_RECORD     = "pushrecord"          // 推送记录
//	DB_TABLE_ANNOUNCEMENT    = "announcement"        // 公告
//	DB_TABLE_PROP            = "propconfig"          // 道具
//	DB_TABLE_LOTTO_SLOTITEMS = "lottoslotitemconfig" // 格子子项
//	DB_TABLE_LOTTO_WEIGHT    = "lottoweightconfig"   // 抽奖权重列表
//	DB_TABLE_LOTTO_STAGE     = "lottostage"          // 抽奖阶段列表
//	//DB_TABLE_LOTTO_INFO               = "lottoinfo"              // 用户抽奖信息
//	//DB_TABLE_LOTTO_SYS_INFO           = "syslottoinfo"           // 系统抽奖信息
//	DB_TABLE_LOTTO_TRANSACTION = "lottotransaction" // 抽奖事务
//	//DB_TABLE_NOTICE                   = "notice"                 // 通知
//	//DB_TABLE_MAIL                     = "mail"                   // 邮件
//	DB_TABLE_RUNE_CONFIG              = "runeconfig"             // 符文配置
//	DB_TABLE_BEFOREGAME_RANDOM_WEIGHT = "beforegamerandomweight" // 赛前随机道具权重配置信息
//	DB_TABLE_PICKUP_WEIGHT            = "pickupweight"           // 游戏中收集物权重配置信息

//	//角色列表信息相关
//	DB_TABLE_ROLEINFO         = "roleinfo"       //角色列表(动态)
//	DB_TABLE_ROLEINFO_CONFIG  = "roleinfoconfig" //角色基础信息（静态）
//	DB_TABLE_ROLE_LEVEL_BONUS = "rolelevelbonus" //角色加成信息（静态）

//	//系统邮件信息相关
//	//DB_TABLE_MAILINFO        = "mailinfo"       //系统邮件列表(动态)
//	DB_TABLE_MAILINFO_CONFIG = "mailinfoconfig" //系统邮件基础信息（静态）
//	DB_TABLE_SYSTEMMAIL      = "systemmaillist" //系统邮件基础信息

//	//好友邮件信息相关（体力赠送）
//	DB_TABLE_STAMINA_GIVEAPPLY_LOG = "staminagiveapplylog" //体力赠送行为日志（动态）
//	DB_TABLE_FRIENDMAIL            = "friendmail"          //好友邮件列表（动态）

//	//签到活动配置信息
//	DB_TABLE_SIGN_IN_ACTIVITY = "signinactivityconfig" // 活动表
//	DB_TABLE_SIGN_IN_ITEM     = "signawardconfig"      // 活动奖励表

//)
//const (
//	//DB_TABLE_GIFT_LOG     = "giftlog"     // 体力请求日志
//	DB_TABLE_GAME_LOG     = "gamelog"     // 游戏日志
//	DB_TABLE_SHOPPING_LOG = "shoppinglog" // 购买日志
//	DB_TABLE_ACCOUNT_LOG  = "accountlog"  // 账户日志
//	//DB_TABLE_DIAMOND_LOG  = "diamondlog"
//	DB_TABLE_IAP_LOG = "iaplog"
//)
const (
	SERVER_ID = 1
)
