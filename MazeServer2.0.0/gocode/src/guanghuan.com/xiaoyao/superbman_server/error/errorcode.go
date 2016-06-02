// errorcode
package xyerror

//错误码分段（以20为一个区段进行分配）
const (
	IAP_STEP_VELIDATE_DATA          = 0   //内购(数据阶段)
	IAP_STEP_VERIFY_RECEIPT_PREPARE = 100 //内购(准备阶段)
	IAP_STEP_VERIFY_RECEIPT_POST    = 200 //内购(后置处理阶段)
	LOTTO_ERR_BASE                  = 300 //抽奖
	GIFT_ERR_BASE                   = 400 //体力
	GAME_ERR_BASE                   = 500 //游戏
	ACCOUNT_ERR_BASE                = 600 //账户
	USR_PROP_ERR_BASE               = 700 //用户道具
	ROLE_LIST_ERR_BASE              = 800 //角色列表
	DAILY_AWARD_ERR_BASE            = 900 //每日登录奖励
)

//体力相关错误码
const (
	GIFT_SUCCESS           = 0
	GIFT_NOT_ALLOWED       = GIFT_ERR_BASE + 0
	GIFT_NOT_EXIST         = GIFT_ERR_BASE + 1
	GIFT_EXPIRED           = GIFT_ERR_BASE + 2
	GIFT_SERVER_ERROR      = GIFT_ERR_BASE + 3
	GIFT_USER_NOT_EXIST    = GIFT_ERR_BASE + 4
	GIFT_ALREADY_REQUESTED = GIFT_ERR_BASE + 5
)

//游戏相关错误码
const (
	GAME_SUCCESS                 = 0
	GAME_NOT_ALLOWED             = GAME_ERR_BASE + 0
	GAME_NOT_EXIST               = GAME_ERR_BASE + 1
	GAME_NOT_ENOUGH_STAMINA      = GAME_ERR_BASE + 2 // 体力不足
	GAME_ALREADY_UPDATED         = GAME_ERR_BASE + 3
	GAME_INVALID_DATA            = GAME_ERR_BASE + 4
	GAME_USER_NOT_EXIST          = GAME_ERR_BASE + 5
	GAME_UPDATE_ACCOUNT_FAIL     = GAME_ERR_BASE + 6
	GAME_UPDATE_GAME_FAIL        = GAME_ERR_BASE + 7
	GAME_NEW_GAME_FAIL           = GAME_ERR_BASE + 8
	GAME_UPDATE_TRANSACTION_FAIL = GAME_ERR_BASE + 9  //下发到transaction服务失败
	GAME_BAD_INPUT_DATA          = GAME_ERR_BASE + 10 //请求数据异常
)

//账户相关错误码
const (
	ACCOUNT_SUCCESS                  = 0
	ACCOUNT_NOT_ALLOWED              = ACCOUNT_ERR_BASE + 0
	ACCOUNT_INVALID_DATA             = ACCOUNT_ERR_BASE + 1
	ACCOUNT_FAIL_TO_GET_UID          = ACCOUNT_ERR_BASE + 2
	ACCOUNT_FAIL_TO_GET_ACCOUNT      = ACCOUNT_ERR_BASE + 3
	ACCOUNT_FAIL_TO_ADD              = ACCOUNT_ERR_BASE + 4
	ACCOUNT_DB_ERROR                 = ACCOUNT_ERR_BASE + 5
	ACCOUNT_VERSION_NOT_SUPPORT      = ACCOUNT_ERR_BASE + 6
	ACCOUNT_FAIL_SEND_TO_TRANSACTION = ACCOUNT_ERR_BASE + 7 //请求transaction服务失败
)

//内购相关错误码
const (
	IAP_SUCCESS                     = 0
	IAP_DB_ERROR                    = IAP_STEP_VELIDATE_DATA + 1
	IAP_INVALID_USER                = IAP_STEP_VELIDATE_DATA + 2
	IAP_INVALID_TRANSACTIONID       = IAP_STEP_VELIDATE_DATA + 3 //请求的transaction_id非法，已经提交过
	IAP_INVALID_GOODS               = IAP_STEP_VELIDATE_DATA + 4
	IAP_INVALID_RECEIPT             = IAP_STEP_VELIDATE_DATA + 5
	IAP_INVALID_IAPVALIDATENODE     = IAP_STEP_VELIDATE_DATA + 6
	IAP_PB_FAILED                   = IAP_STEP_VELIDATE_DATA + 7
	IAP_REQUEST_TRANSACTION_TIMEOUT = IAP_STEP_VELIDATE_DATA + 8 //请求transaction服务超时

	IAP_FAILED_BEFORE_INVOKE_APPLE_VARIFY = IAP_STEP_VERIFY_RECEIPT_PREPARE + 0
	IAP_FAILED_TO_INVOKE_APPLE_VERIFY     = IAP_STEP_VERIFY_RECEIPT_PREPARE + 1
	IAP_APPLE_VARIFY_INVALID_RESPONSE     = IAP_STEP_VERIFY_RECEIPT_PREPARE + 2

	IAP_APPLE_VARIFY_FAIL               = IAP_STEP_VERIFY_RECEIPT_POST + 0
	IAP_APPLE_VARIFY_BUNDLEID_FAIL      = IAP_STEP_VERIFY_RECEIPT_POST + 1 //bundle_id校验失败
	IAP_APPLE_VARIFY_PRODUCTID_FAIL     = IAP_STEP_VERIFY_RECEIPT_POST + 2 //product_id校验失败
	IAP_APPLE_VARIFY_TRANSACTIONID_FAIL = IAP_STEP_VERIFY_RECEIPT_POST + 3 //transaction_id校验失败

	IAP_UPDATE_ACCOUNT_FAIL = IAP_STEP_VERIFY_RECEIPT_POST + 5
)

//抽奖相关错误码
const (
	LOTTO_SUCCESS                         = 0
	LOTTO_INVALID_NODE                    = LOTTO_ERR_BASE + 0  //lotto事务节点名称获取失败
	LOTTO_INVALID_USER                    = LOTTO_ERR_BASE + 1  //lotto请求uid不合法
	LOTTO_REQUEST_TRANSACTION_TIMEOUT     = LOTTO_ERR_BASE + 2  //lotto请求transaction服务超时
	LOTTO_PB_FAILED                       = LOTTO_ERR_BASE + 3  //lotto请求pb编解码失败
	LOTTO_UNKOWN_CMD                      = LOTTO_ERR_BASE + 4  //lotto请求未知的cmd
	LOTTO_UNKOWN_RES_OP_TYPE              = LOTTO_ERR_BASE + 5  //lotto未知的资源请求类型
	LOTTO_DB_ERR_RES_OP_TYPE              = LOTTO_ERR_BASE + 6  //lotto数据库操作错误
	LOTTO_DB_ERR_LOAD_PROPS               = LOTTO_ERR_BASE + 7  //lotto prop数据加载失败
	LOTTO_DB_ERR_LOAD_SLOTITEMS           = LOTTO_ERR_BASE + 8  //lotto slotitems数据加载失败
	LOTTO_DB_ERR_LOAD_WEIGHT              = LOTTO_ERR_BASE + 9  //lotto weight数据加载失败
	LOTTO_DB_ERR_DrawAwardType            = LOTTO_ERR_BASE + 10 //lotto 错误的DrawAwardType
	LOTTO_DB_ERR_SYS_LOTTO_INFO_TIMESTAMP = LOTTO_ERR_BASE + 11 //lotto 查询syslottoinfo的时间戳信息失败
	LOTTO_DB_ERR_SYS_LOTTO_INFO           = LOTTO_ERR_BASE + 12 //lotto syslottoinfo数据库操作失败
	LOTTO_DB_ERR_LOTTO_TRANSACTION_LOG    = LOTTO_ERR_BASE + 13 //lotto 记录抽奖事务错误
	LOTTO_DB_ERR_LOTTO_LOG                = LOTTO_ERR_BASE + 14 //lotto 记录抽奖日志
)

//角色列表相关错误码
const (
	ROLE_LIST_SUCCESS      = 0
	ROLE_LIST_INVALID_NODE = ROLE_LIST_ERR_BASE + 0     //rolelist事务节点名称获取失败
	ROLE_LIST_INVALID_USER = ROLE_LIST_INVALID_NODE + 1 //rolelist请求uid不合法
	ROLE_LIST_ERR_ROLE_ID  = ROLE_LIST_INVALID_USER + 1 //roleid 不对
	ROLE_LIST_ERR_CMD      = ROLE_LIST_ERR_ROLE_ID + 1  //cmd指令错误
)

const (
	USR_SUCCESS        = 0
	USR_PROP_ERR_CACHE = USR_PROP_ERR_BASE + 1 //prop cache错误
)

const (
	MISSION_SUCCESS                                       = 0
	MISSION_DAILY_AWARD_ERR_STATIC_INFO                   = DAILY_AWARD_ERR_BASE + 0 //查询每日登录奖励静态信息失败
	MISSION_DAILY_AWARD_ERR_QUERY_USER_DAILY_LOGIN_RECORD = DAILY_AWARD_ERR_BASE + 1 //查询玩家每日登录信息失败
	MISSION_DAILY_AWARD_ERR_SET_USER_LOGIN_BITMAP         = DAILY_AWARD_ERR_BASE + 2 //设置玩家领奖位失败
	MISSION_DAILY_AWARD_ERR_COLLECT_INVALID_REQUEST       = DAILY_AWARD_ERR_BASE + 3 //每日登录领奖，请求消息错误
	MISSION_USER_MISSION_ERR_QUERY                        = DAILY_AWARD_ERR_BASE + 4 //查询玩家任务信息失败
)
