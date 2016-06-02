// config
package main

import (
	//proto "code.google.com/p/goprotobuf/proto"
	"flag"
	"fmt"
	//httplib "github.com/astaxie/beego/httplib"
	//battery "guanghuan.com/xiaoyao/superbman_server/battery_run_net"
	//"strings"
	beegoconf "github.com/astaxie/beego/config"
	"guanghuan.com/xiaoyao/common/log"
	"guanghuan.com/xiaoyao/superbman_server/error"
	"os"
)

type Config struct {
	DBConfig
	Url, IniFile, DBUrl, DBName, ApnsNatsUrl string
	Id, Sum                                  uint64
	Cmd                                      int
	FlowControl                              int
	ConnTimeout, RwTimeout                   int64
}

var DefConfig = Config{
	Url: "192.168.1.205:10003",
}

const (
	API_URI_LOGIN                   = "/v1/login/:token"
	API_URI_GET_USER_GAMEDATA       = "/v1/user/:token"
	API_URI_GET_FRIEND_GAMEDATA     = "/v1/friend/:token"
	API_URI_NEWGAME                 = "/v1/newgame/:token"
	API_URI_ADD_GAMEDATA            = "/v1/gameresult/:token"
	API_URI_ADD_GAMEDATA2           = "/v2/gameresult/:token"
	API_URI_STAMINA                 = "/v1/stamina/:token"
	API_URI_GIFT_QUERY              = "/v1/gift/query/:token"   // 查询
	API_URI_GIFT_OP                 = "/v1/gift/op/:token"      // 确认
	API_URI_GOODS_QUERY             = "/v1/goods/query/:token"  // 查询商品列表
	API_URI_GOODS_BUY               = "/v1/goods/buy/:token"    // 购买商品
	API_URI_ANNOUNCEMENT            = "/v2/announcement/:token" // 通告
	API_URI_DEVICE_ID_SUBMIT        = "/v1/device/device_id"
	API_URI_IAP_VERIFY              = "/iap_verify/order_verify"             //iap内购校验
	API_URI_LOTTO_OP                = "/v2/lotto/lotto_op"                   //抽奖请求
	API_URI_LOTTO_RES_OP            = "/v2/lotto/lotto_res_op"               //抽奖资源请求
	API_URI_QUERY_DAILY_AWARD       = "/v2/dailyaward/query"                 //查询玩家每日登录奖励
	API_URI_QUERY_USER_MISSION      = "/v2/usermission/query"                //查询玩家任务列表
	API_URI_CONFIRM_USER_MISSION    = "/v2/usermission/confirm"              //领取玩家任务奖励
	API_URI_PROP_RES_QUERY          = "/v2/prop/prop_res_query"              //道具信息查询
	API_URI_SIGNIN_ACTIVITY_QUERY   = "/v2/signin/query"                     //玩家签到活动查询
	API_URI_SIGNIN                  = "/v2/signin/sign"                      //玩家签到
	API_URI_CHECKPOINT_QUERY_RANGE  = "/v2/checkpoint/query_range"           //查询玩家区间记忆点信息
	API_URI_CHECKPOINT_QUERY_DETAIL = "/v2/checkpoint/query_detail"          //查询记忆点排行榜
	API_URI_CHECKPOINT_COMMIT       = "/v2/checkpoint/commit"                //提交记忆点数据
	API_URI_WALLET_QUERY            = "/v2/wallet/query"                     //查询玩家钱包数据
	API_URI_BEFOREGAME_OP           = "/v2/beforegameprop/beforegameprop_op" //游戏前商品相关信息
	API_URI_ROLE_INFO               = "/v2/roleinfolist/roleinfolist_op"     //玩家的角色信息相关请求
	API_URI_JIGSAW                  = "/v2/jigsaw/jigsaw_op"                 //玩家的拼图信息相关请求
	API_URI_FRIEND_MAIL             = "/v2/friendmail/friendmail_op"         //玩家的好友邮件相关请求
	API_URI_SYS_MAIL                = "/v2/systemmail/systemmail_op"         //玩家的系统邮件相关请求
	API_URI_RUNE                    = "/v2/rune/rune_op"                     //玩家的符文相关请求
	API_URI_MEMCACHE                = "/v2/memcache/memcache_op"             //玩家的memcache相关请求
	API_URI_ADVERTISEMENT           = "/v2/advertisement/advertisement"      //广告请求
	API_URI_BIND                    = "/v2/bind/bind"                        //广告请求

	API_URI_MAINTENANCE = "/maintenance" //运营接口。修改道具数。
)

var MAPCMD2DESCRIPTION = map[int]string{
	TEST_CMD_CODE_LOGIN:                    "login",
	TEST_CMD_CODE_PPROF_LOGIN:              "pprof login",
	TEST_CMD_CODE_BIND:                     "bind",
	TEST_CMD_CODE_PPROF_BIND:               "pprof bind",
	TEST_CMD_CODE_QUERY_ANNOUNCEMENT:       "query announcement",
	TEST_CMD_CODE_PPROF_QUERY_ANNOUNCEMENT: "pprof query announcement",

	TEST_CMD_CODE_QUERY_BEFOREGAMEGOODS:       "query beforegamegoods",
	TEST_CMD_CODE_PPROF_QUERY_BEFOREGAMEGOODS: "pprof query beforegamegoods",
	TEST_CMD_CODE_BUY_BEFOREGAMEGOODS:         "buy beforegamegoods",
	TEST_CMD_CODE_PPROF_BUY_BEFOREGAMEGOODS:   "pprof buy beforegamegoods",
	TEST_CMD_CODE_USE_BEFOREGAMEGOODS:         "use beforegamegoods",
	TEST_CMD_CODE_PPROF_USE_BEFOREGAMEGOODS:   "pprof use beforegamegoods",

	TEST_CMD_CODE_QUERY_FRIENDMAIL:                     "query friend mail",
	TEST_CMD_CODE_PPROF_QUERY_FRIENDMAIL:               "pprof query friend mail",
	TEST_CMD_CODE_GIVE_FRIEND_FROM_FRIENDSHIP:          "give friend from friendship",
	TEST_CMD_CODE_PPROF_GIVE_FRIEND_FROM_FRIENDSHIP:    "pprof give friend from friendship",
	TEST_CMD_CODE_APPLY_FRIEND_FROM_FRIENDSHIP:         "apply friend from friendship",
	TEST_CMD_CODE_PPROF_APPLY_FRIEND_FROM_FRIENDSHIP:   "pprof apply friend from friendship",
	TEST_CMD_CODE_GIVE_FRIEND_FROM_FRIENDMAIL:          "give friend from friendmail",
	TEST_CMD_CODE_PPROF_GIVE_FRIEND_FROM_FRIENDMAIL:    "pprof give friend from friendmail",
	TEST_CMD_CODE_APPLY_FRIEND_FROM_FRIENDMAIL:         "apply friend from friendmail",
	TEST_CMD_CODE_PPROF_APPLY_FRIEND_FROM_FRIENDMAIL:   "pprof apply friend from friendmail",
	TEST_CMD_CODE_COLLECT_FRIEND_FROM_FRIENDMAIL:       "collect friend from friendmail",
	TEST_CMD_CODE_PPROF_COLLECT_FRIEND_FROM_FRIENDMAIL: "pprof collect friend from friendmail",
	//TEST_CMD_CODE_CONFIRM_ALLFRIENDMAIL:                "confirm all friendmail",
	//TEST_CMD_CODE_PPROF_CONFIRM_ALLFRIENDMAIL:          "pprof confirm all friendmail",
	//iap
	TEST_CMD_CODE_IAP_VARIFY:       "iap varify",
	TEST_CMD_CODE_PPROF_IAP_VARIFY: "pprof iap varify",
	//goods
	TEST_CMD_CODE_QUERY_GOODS:       "query goods",
	TEST_CMD_CODE_PPROF_QUERY_GOODS: "pprof query goods",
	TEST_CMD_CODE_BUY_GOODS:         "buy goods",
	TEST_CMD_CODE_PPROF_BUY_GOODS:   "pprof buy goods",
	//stamina
	TEST_CMD_CODE_QUERY_STAMINA: "query stamina",
	//frienddata
	TEST_CMD_CODE_QUERY_FRIENDDATA:       "query frienddata",
	TEST_CMD_CODE_PPROF_QUERY_FRIENDDATA: "pprof query frienddata",
	//syslotto
	TEST_CMD_CODE_SYSLOTTO:       "syslotto",
	TEST_CMD_CODE_PPROF_SYSLOTTO: "pprof syslotto",
	//game
	TEST_CMD_CODE_NEWGAME:       "new game",
	TEST_CMD_CODE_PPROF_NEWGAME: "pprof new game",
	//TEST_CMD_CODE_NEWGAMEANDRESULT:       "new game and result",
	TEST_CMD_CODE_PPROF_NEWGAMEANDRESULT: "pprof new game and  result",
	TEST_CMD_CODE_GAMERESULT:             "game result",
	//signin
	TEST_CMD_CODE_QUERY_SIGNIN_ACTIVITY:       "query signin activity",
	TEST_CMD_CODE_PPROF_QUERY_SIGNIN_ACTIVITY: "pprof query signin activity",
	TEST_CMD_CODE_SIGNIN_ACTIVITY:             "signin activity",
	TEST_CMD_CODE_PPROF_SIGNIN_ACTIVITY:       "pprof signin activity",
	//usermission
	TEST_CMD_CODE_QUERY_USERMISSION:         "query user mission",
	TEST_CMD_CODE_PPROF_QUERY_USERMISSION:   "pprof query user mission",
	TEST_CMD_CODE_CONFIRM_USERMISSION:       "confirm user mission",
	TEST_CMD_CODE_PPROF_CONFIRM_USERMISSION: "pprof confirm user mission",
	//checkpoint
	TEST_CMD_CODE_QUERY_USERCHECKPOINT:               "query user checkpoint",
	TEST_CMD_CODE_PPROF_QUERY_USERCHECKPOINT:         "pprof query user checkpoint",
	TEST_CMD_CODE_QUERY_CHECKPOINT_FRINENDRANK:       "query checkpoint friend rank",
	TEST_CMD_CODE_PPROF_QUERY_CHECKPOINT_FRINENDRANK: "pprof query checkpoint friend rank",
	TEST_CMD_CODE_QUERY_CHECKPOINT_GLOBALRANK:        "query checkpoint global rank",
	TEST_CMD_CODE_PPROF_QUERY_CHECKPOINT_GLOBALRANK:  "pprof query checkpoint global rank",
	//wallet
	TEST_CMD_CODE_QUERY_WALLET:       "query wallet",
	TEST_CMD_CODE_PPROF_QUERY_WALLET: "pprof query wallet",
	//roleinfo
	TEST_CMD_CODE_QUERY_USER_ROLEINFO:              "query user roleinfo",
	TEST_CMD_CODE_PPROF_QUERY_USER_ROLEINFO:        "pprof query user roleinfo",
	TEST_CMD_CODE_QUERY_FRIEND_ROLEINFO:            "query friend roleinfo",
	TEST_CMD_CODE_PPROF_QUERY_FRIEND_ROLEINFO:      "pprof query friend roleinfo",
	TEST_CMD_CODE_SET_USER_SELECTED_ROLEINFO:       "set user selected roleinfo",
	TEST_CMD_CODE_PPROF_SET_USER_SELECTED_ROLEINFO: "pprof set user selected roleinfo",
	TEST_CMD_CODE_BUY_ROLE:                         "buy rune",
	TEST_CMD_CODE_PPROF_BUY_ROLE:                   "pprof buy rune",
	//jigsaw
	TEST_CMD_CODE_QUERY_USER_JIGSAW:       "query user jigsaw",
	TEST_CMD_CODE_PPROF_QUERY_USER_JIGSAW: "pprof query user jigsaw",
	TEST_CMD_CODE_BUY_JIGSAW:              "buy jigsaw",
	TEST_CMD_CODE_PPROF_BUY_JIGSAW:        "pprof buy jigsaw",
	//sysmail
	TEST_CMD_CODE_QUERY_USER_SYSMAIL:         "query user sysmail",
	TEST_CMD_CODE_PPROF_QUERY_USER_SYSMAIL:   "pprof query user sysmail",
	TEST_CMD_CODE_CONFIRM_USER_SYSMAIL:       "confirm user sysmail",
	TEST_CMD_CODE_PPROF_CONFIRM_USER_SYSMAIL: "pprof confirm user sysmail",
	TEST_CMD_CODE_READ_USER_SYSMAIL:          "read user sysmail",
	TEST_CMD_CODE_PPROF_READ_USER_SYSMAIL:    "pprof read user sysmail",
	//rune
	TEST_CMD_CODE_QUERY_RUNE:       "query rune",
	TEST_CMD_CODE_PPROF_QUERY_RUNE: "pprof query rune",
	TEST_CMD_CODE_BUY_RUNE:         "buy rune",
	TEST_CMD_CODE_PPROF_BUY_RUNE:   "pprof buy rune",

	//memcache
	TEST_CMD_CODE_MEMCACHE_GET:        "memcache get",
	TEST_CMD_CODE_PPROF_MEMCACHE_GET:  "pprof memcache get",
	TEST_CMD_CODE_MEMCACHE_SET:        "memcache set",
	TEST_CMD_CODE_PPROF_MEMCACHE_SET:  "pprof memcache set",
	TEST_CMD_CODE_MEMCACHES_GET:       "memcaches get",
	TEST_CMD_CODE_PPROF_MEMCACHES_GET: "pprof memcaches get",
	TEST_CMD_CODE_MEMCACHES_SET:       "memcaches set",
	TEST_CMD_CODE_PPROF_MEMCACHES_SET: "pprof memcaches set",

	//db
	TEST_CMD_CODE_DB_QUERY_USER_ACCOUNT_SOMEFIELD: "query user account somefields",

	//maintenance
	TEST_CMD_CODE_MAINTENANCE_PROP: "maintenance prop",

	TEST_CMD_CODE_APN_NOTIFICATION:            "apn notification",
	TEST_CMD_CODE_APN_NOTIFICATION2APNS:       "apn notification to apns",
	TEST_CMD_CODE_APN_ENABLEDEVICETOKEN2APNS:  "enable devicetoken to apns",
	TEST_CMD_CODE_APN_DISABLEDEVICETOKEN2APNS: "disable devicetoken to apns",

	TEST_CMD_CODE_TIMER: "timer",

	TEST_CMD_CODE_ADVERTISEMENT: "advertisement",

	TEST_CMD_CODE_IAPSTATISTIC: "iapstatistic",

	TEST_CMD_CODE_DB: "db",
}

func CmdDescription() (str string) {

	for i := 0; i < TEST_CMD_CODE_MAX; i++ {
		if v, ok := MAPCMD2DESCRIPTION[i]; ok {
			str += fmt.Sprintf("%d %s\n", i, v)
		} else {
			str += fmt.Sprintf("%d unkown\n", i)
		}
	}
	return
}

const (
	//login
	TEST_CMD_CODE_LOGIN = iota
	TEST_CMD_CODE_PPROF_LOGIN
	TEST_CMD_CODE_BIND
	TEST_CMD_CODE_PPROF_BIND
	//annoucement
	TEST_CMD_CODE_QUERY_ANNOUNCEMENT
	TEST_CMD_CODE_PPROF_QUERY_ANNOUNCEMENT
	//beforegamegoods
	TEST_CMD_CODE_QUERY_BEFOREGAMEGOODS
	TEST_CMD_CODE_PPROF_QUERY_BEFOREGAMEGOODS
	TEST_CMD_CODE_BUY_BEFOREGAMEGOODS
	TEST_CMD_CODE_PPROF_BUY_BEFOREGAMEGOODS
	TEST_CMD_CODE_USE_BEFOREGAMEGOODS
	TEST_CMD_CODE_PPROF_USE_BEFOREGAMEGOODS
	//friend mail
	TEST_CMD_CODE_QUERY_FRIENDMAIL
	TEST_CMD_CODE_PPROF_QUERY_FRIENDMAIL
	TEST_CMD_CODE_GIVE_FRIEND_FROM_FRIENDSHIP
	TEST_CMD_CODE_PPROF_GIVE_FRIEND_FROM_FRIENDSHIP
	TEST_CMD_CODE_APPLY_FRIEND_FROM_FRIENDSHIP
	TEST_CMD_CODE_PPROF_APPLY_FRIEND_FROM_FRIENDSHIP
	TEST_CMD_CODE_GIVE_FRIEND_FROM_FRIENDMAIL
	TEST_CMD_CODE_PPROF_GIVE_FRIEND_FROM_FRIENDMAIL
	TEST_CMD_CODE_APPLY_FRIEND_FROM_FRIENDMAIL
	TEST_CMD_CODE_PPROF_APPLY_FRIEND_FROM_FRIENDMAIL
	TEST_CMD_CODE_COLLECT_FRIEND_FROM_FRIENDMAIL
	TEST_CMD_CODE_PPROF_COLLECT_FRIEND_FROM_FRIENDMAIL
	//TEST_CMD_CODE_CONFIRM_ALLFRIENDMAIL
	//TEST_CMD_CODE_PPROF_CONFIRM_ALLFRIENDMAIL
	//iap
	TEST_CMD_CODE_IAP_VARIFY
	TEST_CMD_CODE_PPROF_IAP_VARIFY
	//goods
	TEST_CMD_CODE_QUERY_GOODS
	TEST_CMD_CODE_PPROF_QUERY_GOODS
	TEST_CMD_CODE_BUY_GOODS
	TEST_CMD_CODE_PPROF_BUY_GOODS
	//stamina
	TEST_CMD_CODE_QUERY_STAMINA
	//frienddata
	TEST_CMD_CODE_QUERY_FRIENDDATA
	TEST_CMD_CODE_PPROF_QUERY_FRIENDDATA
	//syslotto
	TEST_CMD_CODE_SYSLOTTO
	TEST_CMD_CODE_PPROF_SYSLOTTO
	//game
	TEST_CMD_CODE_NEWGAME
	TEST_CMD_CODE_PPROF_NEWGAME
	TEST_CMD_CODE_GAMERESULT
	//TEST_CMD_CODE_NEWGAMEANDRESULT
	TEST_CMD_CODE_PPROF_NEWGAMEANDRESULT

	//signin
	TEST_CMD_CODE_QUERY_SIGNIN_ACTIVITY
	TEST_CMD_CODE_PPROF_QUERY_SIGNIN_ACTIVITY
	TEST_CMD_CODE_SIGNIN_ACTIVITY
	TEST_CMD_CODE_PPROF_SIGNIN_ACTIVITY
	//usermission
	TEST_CMD_CODE_QUERY_USERMISSION
	TEST_CMD_CODE_PPROF_QUERY_USERMISSION
	TEST_CMD_CODE_CONFIRM_USERMISSION
	TEST_CMD_CODE_PPROF_CONFIRM_USERMISSION
	//checkpoint
	TEST_CMD_CODE_QUERY_USERCHECKPOINT
	TEST_CMD_CODE_PPROF_QUERY_USERCHECKPOINT
	TEST_CMD_CODE_QUERY_CHECKPOINT_FRINENDRANK
	TEST_CMD_CODE_PPROF_QUERY_CHECKPOINT_FRINENDRANK
	TEST_CMD_CODE_QUERY_CHECKPOINT_GLOBALRANK
	TEST_CMD_CODE_PPROF_QUERY_CHECKPOINT_GLOBALRANK
	//wallet
	TEST_CMD_CODE_QUERY_WALLET
	TEST_CMD_CODE_PPROF_QUERY_WALLET
	//roleinfo
	TEST_CMD_CODE_QUERY_USER_ROLEINFO
	TEST_CMD_CODE_PPROF_QUERY_USER_ROLEINFO
	TEST_CMD_CODE_QUERY_FRIEND_ROLEINFO
	TEST_CMD_CODE_PPROF_QUERY_FRIEND_ROLEINFO
	TEST_CMD_CODE_SET_USER_SELECTED_ROLEINFO
	TEST_CMD_CODE_PPROF_SET_USER_SELECTED_ROLEINFO
	TEST_CMD_CODE_BUY_ROLE
	TEST_CMD_CODE_PPROF_BUY_ROLE
	//jigsaw
	TEST_CMD_CODE_QUERY_USER_JIGSAW
	TEST_CMD_CODE_PPROF_QUERY_USER_JIGSAW
	TEST_CMD_CODE_BUY_JIGSAW
	TEST_CMD_CODE_PPROF_BUY_JIGSAW
	//sysmail
	TEST_CMD_CODE_QUERY_USER_SYSMAIL
	TEST_CMD_CODE_PPROF_QUERY_USER_SYSMAIL
	TEST_CMD_CODE_CONFIRM_USER_SYSMAIL
	TEST_CMD_CODE_PPROF_CONFIRM_USER_SYSMAIL
	TEST_CMD_CODE_READ_USER_SYSMAIL
	TEST_CMD_CODE_PPROF_READ_USER_SYSMAIL
	//rune
	TEST_CMD_CODE_QUERY_RUNE
	TEST_CMD_CODE_PPROF_QUERY_RUNE
	TEST_CMD_CODE_BUY_RUNE
	TEST_CMD_CODE_PPROF_BUY_RUNE

	//memcache
	TEST_CMD_CODE_MEMCACHE_GET
	TEST_CMD_CODE_PPROF_MEMCACHE_GET
	TEST_CMD_CODE_MEMCACHE_SET
	TEST_CMD_CODE_PPROF_MEMCACHE_SET
	TEST_CMD_CODE_MEMCACHES_GET
	TEST_CMD_CODE_PPROF_MEMCACHES_GET
	TEST_CMD_CODE_MEMCACHES_SET
	TEST_CMD_CODE_PPROF_MEMCACHES_SET
	//alert
	TEST_CMD_CODE_SEND_ALERT

	//db
	TEST_CMD_CODE_DB_QUERY_USER_ACCOUNT_SOMEFIELD

	//maintenance
	TEST_CMD_CODE_MAINTENANCE_PROP

	//apn notification
	TEST_CMD_CODE_APN_NOTIFICATION
	TEST_CMD_CODE_APN_NOTIFICATION2APNS
	TEST_CMD_CODE_APN_ENABLEDEVICETOKEN2APNS
	TEST_CMD_CODE_APN_DISABLEDEVICETOKEN2APNS

	//timer
	TEST_CMD_CODE_TIMER

	//advertisement
	TEST_CMD_CODE_ADVERTISEMENT

	//iapstatistic
	TEST_CMD_CODE_IAPSTATISTIC

	//db
	TEST_CMD_CODE_DB

	TEST_CMD_CODE_MAX
)

func initConfig() {
	ProcessCmd()
	ParseConfigFile()
	xylog.ApplyConfig(xylog.DefConfig)

}

func ProcessCmd() (should_continue bool) {
	//[caution!]只有各个服务节点不同的配置项才放在进程命令行中，不同服务节点相同的配置项放在配置文件中
	flag.StringVar(&DefConfig.IniFile, "config", os.Args[0]+".ini", "ini file")
	flag.Uint64Var(&DefConfig.Id, "id", 0, "the begin of id")
	flag.Uint64Var(&DefConfig.Sum, "sum", 100000, "count of request")
	flag.IntVar(&DefConfig.FlowControl, "fc", 6000, "flow control")
	flag.IntVar(&DefConfig.Cmd, "cmd", 0, CmdDescription())
	xylog.DefConfig.ProcessCmd()
	flag.Parse()
	should_continue = true
	return
}

func ApplyConfig() {
	xylog.ApplyConfig(xylog.DefConfig)
	xylog.Info("Config: %s", xylog.DefConfig.String())
}

//解析配置文件
//[caution!!!]对于所有服务节点相同的配置项放在配置文件中
func ParseConfigFile() error {
	configs, err := beegoconf.NewConfig("ini", DefConfig.IniFile)
	if err != xyerror.ErrOK {
		fmt.Printf("load config file %s failed : %v\n", DefConfig.IniFile, err)
		return err
	}

	DefConfig.Url = configs.String("Server::url")
	DefConfig.DBUrl = configs.String("Server::dburl")
	DefConfig.DBName = configs.String("Server::dbname")
	DefConfig.ApnsNatsUrl = configs.String("Server::apnsnatsurl")
	if DefConfig.RwTimeout, err = configs.Int64("Server::rwtimeout"); err != nil {
		fmt.Printf("load Server::rwtimeout failed : %v\n", err)
		return err
	}
	if DefConfig.ConnTimeout, err = configs.Int64("Server::conntimeout"); err != nil {
		fmt.Printf("load Server::conntimeout failed : %v\n", err)
		return err
	}

	err = xylog.DefConfig.ProcessIniConfig(configs)
	if err != xyerror.ErrOK {
		fmt.Printf("DefConfig.LogConfig.ProcessIniConfig failed : %v\n", err)
		return err

	}

	fmt.Printf("DefConfig : %v", DefConfig)

	return err
}
