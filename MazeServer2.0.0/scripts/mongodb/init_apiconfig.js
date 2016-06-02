/*
  name 配置项名称
  curclientversion 版本信息
  defaultdiamond 默认钻石数
  defaultcoin 默认金币数
  defaultgiftacceptmaxcount 默认好友邮件数目
  defaultgiftaskcooldown 默认请求体力时间间隔，24hrs
  defaultgiftgivecooldown 默认赠送体力时间间隔，2hrs
  defaultgiftvalidtimesec 体力请求的有效期
  defaultgiftnotifycount 体力请求发送apn通知的阈值
  defaultmaxregstamina 默认体力上限
  defaultstamina 默认体力
  defaultstaminregintervalsec 无用
  defaultsysaccount 默认系统账户
  enablesecurity 是否对通讯数据进行加密
  giftmaxbatchsize
  isproduction 是否是生产环境，内购
  isapnproduction apns是否是生产环境
  loglevel 日志级别
  mailboxlimit 邮箱最多显示邮件数
  maxfriendsrequestcount 每次最多向多少个好友发起请求(包括查询好友数据，请求体力，等)
  maxrequestsize 每个请求数据的最大值(byte)
  minclientversion 支持最低的客户端版本
  maxclientversion 支持最高的客户端版本
  transactionnodecount 事务服务节点数
  channelcount 事务服务节点 channel数
  channelmaxmsg channel 可缓存的最大消息数
  bundleid bundle_id
  natstimeout nats消息超时时间，单位：秒
  lottoslotcount 抽奖格子数
  lottoinituservalue 玩家默认的奖池价值
  lottocostpertime 每次抽奖消耗的内部价值
  lottodeduct 抽奖抽水的内部价值
  syslottofreecount 默认免费抽奖次数
  syslottorefreshtime 奖品池刷新时间
  aftergamelottodeleteslotlimit 游戏后抽奖删除格子次数
  aftergamelottomallitems 游戏后抽奖删除格子对应商品id列表
  missioncountlimit 活跃玩家任务数目
  dailymissioncountlimit 日常任务限制（无用，目前日常任务在代码中限制为1个）
  dailymissionrefreshhour 日常任务刷新时间，8点
  checkpointidnum 记忆点数目
  checkpointglobalrankreloadsecs 记忆点全局玩家排行榜刷新时间间隔
  checkpointglobalranksize 记忆点全局排行榜显示玩家数目
  aftergameawardfactor 游戏后奖励系数，100/10000
  selectedlottoslotfortest 4test 抽奖选中格子， -1表示不指定
  panic 程序crash打印堆栈开关
  quotasneedfinish 需要记忆点完成才记录的任务指标列表
  appkey 应用的appkey
  appid  应用的标识，数据中心分配
  appmastersecret 应用的appmastersecret
  apnnotifydeviceperreq apn推送每个请求中多少个devicetoken
  appsecretkey 数据中心密钥
  datacenterurl 数据中心链接
  */
db.apiconfig.save({
  "name" : "test2.0",
  "curclientversion" : 10000,
  "defaultdiamond" : 100,
  "defaultcoin" : 1000,
  "defaultgiftacceptmaxcount" : 30,
  "defaultgiftaskcooldown" : NumberLong(86400),
  "defaultgiftgivecooldown" : NumberLong(7200),
  "defaultgiftvalidtimesec" : NumberLong(604800),
  "defaultgiftnotifycount" : 5,
  "defaultmaxregstamina" : 5,
  "defaultstamina" : 5,
  "defaultstaminregintervalsec" : NumberLong(600),
  "defaultsysaccount" : "sys",
  "enablesecurity" : true,
  "giftmaxbatchsize" : 100,
  "isproduction" : true,
  "isapnsproduction" : true,
  "loglevel" : 1,
  "mailboxlimit" : 30,
  "maxfriendsrequestcount" : 10,
  "maxrequestsize" : NumberLong(20480),
  "minclientversion" : 10200,
  "maxclientversion" : 90000,
  "transactionnodecount" : 8,
  "channelcount" : 100,
  "channelmaxmsg" : 2500,
  "bundleid" : "com.guanghuan.SuperBMan,com.737.batteryrun,com.737.batteryrun.cn,com.737.batteryrun2.cn",
  "natstimeout" : 10,
  "lottoslotcount" : 8,
  "lottoinituservalue" : 120,
  "lottocostpertime" : 10,
  "lottodeduct" : 20,
  "syslottofreecount" : 1,
  "syslottorefreshtime" : 10800,
  "aftergamelottodeleteslotlimit" : 3,
  "aftergamelottomallitems" : "1;2;3",
  "missioncountlimit" : 3,
  "dailymissioncountlimit" : 3,
  "dailymissionrefreshhour" : 8,
  "checkpointidnum" : 500,
  "checkpointglobalrankreloadsecs" : 120,
  "checkpointglobalranksize" : 5,
  "globalrankreloadsecs" :120,
  "globalranksize" : 20,
  "aftergameawardfactor" : 100,
  "selectedlottoslotfortest" : -1,
  "panic" : true,
  "quotasneedfinish":[4008,4009,4010,4011,5004,5005],
  "appkey":"787f1414d2ca78bc338eaf255f309553",
  "appid":"10004",
  "appmastersecret":"1ygzzohcqzh8qzeqrjmhmojycmp3burk",
  "apnnotifydeviceperreq":100,
  "appsecretkey":"a0e90ad18ce588000586ab6becb43923",
  "datacenterurl":"http://gateway.dc.737.com/index.php",
  "sdkloginurl":"http://sdk.syapi.737.com/sdk/dckp/",
  "invalidScore":300000
});
