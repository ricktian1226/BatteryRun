package xybusiness

//business code definition
const (
    BusinessCode_Unkown                    = iota
    BusinessCode_Login                     //1 登录/注册
    BusinessCode_NewGame                   //2 新游戏
    BusinessCode_GameResult                //3 提交游戏数据
    BusinessCode_Stamina                   //4 体力查询
    BusinessCode_OldStaminaGift            //5 赠送操作
    BusinessCode_GoodsBuy                  //6 购买商品
    BusinessCode_IapValidate               //7 Iap内购校验
    BusinessCode_Lotto                     //8 抽奖
    BusinessCode_QuerySignIn               //9 查询玩家签到活动
    BusinessCode_SignIn                    //10 玩家签到
    BusinessCode_QueryUserMission          //11 查询玩家任务
    BusinessCode_ConfirmUserMission        //12 领取玩家任务奖励
    BusinessCode_RoleInfoList              //13 角色消息
    BusinessCode_FriendMailInfoList        //14 好友邮件消息
    BusinessCode_Announcement              //15 公告
    BusinessCode_SystemMailInfoList        //16 查询系统邮件
    BusinessCode_QueryGoods                //17 查询商品信息
    BusinessCode_Jigsaw                    //18 拼图
    BusinessCode_Rune                      //19 符文
    BusinessCode_QueryUserCheckPoints      //20 查询玩家区间记忆点信息
    BusinessCode_QueryUserCheckPointDetail //21 查询玩家记忆点详细信息
    BusinessCode_CommitCheckPoint          //22 提交玩家记忆点成绩
    BusinessCode_QueryWallet               //23 查询玩家钱包信息
    BusinessCode_BeforeGameProp            //24 赛前道具
    BusinessCode_GameResult2               //25 提交游戏数据2
    BusinessCode_MemCache                  //26 memcache操作
    BusinessCode_MaintenanceProp           //27 运营操作-物品
    BusinessCode_Bind                      //28 绑定账户

    BusinessCode_QuerySignIn2     //29 新版本查询玩家签到活动
    BusinessCode_ShareQuery       // 30分享信息查询
    BusinessCode_ShareRequest     // 31分享请求
    BusinessCode_CheckPointUnlock // 关卡解锁
    BusinessCode_SDKOrderOp       // sdk订单操作
    BusinessCode_SDKOrderQuery    // sdk订单查询
    BusinessCode_SDKAddOrder      // sdk添加订单
    BusinessCode_GlobalRankList   // 全局排行榜
    BusinessCode_CreatName        // 玩家起名
)
