// tpidmap 玩家第三方账户信息
db.tpidmap.ensureIndex({sid:1,source:1},{background:1});//根据第三方id查询玩家对应tpid信息
db.tpidmap.ensureIndex({gid:1,source:1},{background:1});//根据uid查询玩家对应tpid信息
// useraccomplishment 玩家成就信息
db.useraccomplishment.ensureIndex({uid:1},{background:1});//定时任务加载公告配置信息
// usercheckpoint 玩家记忆点信息
db.usercheckpoint.ensureIndex({uid:1,checkpointid:1},{background:1});//查询玩家记忆点信息列表
db.usercheckpoint.ensureIndex({checkpointid:1, score:-1},{background:1});//查询玩家记忆点信息列表
// memcache 玩家缓存信息
db.memcache.ensureIndex({uid:1, key:1, platform:1},{background:1});
