/*
   mongodb的固定集合是天生为日志数据而生的，所有的日志数据库集合都采用固定集合方式创建。
   日志集合包括：
   accountlog   登录日志
   gamelog      游戏日志
   shoppinglog  购买日志
   iaplog       内购日志
   lottolog     抽奖日志
   maintenancelog 运营日志  
   collection的size和max分析请参照《User数据模型.xls》文档：
   https://192.168.1.204:8443/svn/Savage/Program/MazeServer2.0.0/scripts/mongodb
*/
db.createCollection("accountlog",{size:2147483648, capped:true, max:3000000,autoIndexId:false});
db.createCollection("gamelog",{size:2147483648, capped:true, max:5000000,autoIndexId:false});
db.createCollection("shoppinglog",{size:524288000, capped:true, max:3000000,autoIndexId:false});
db.createCollection("iaplog",{size:524288000, capped:true, max:3000000,autoIndexId:false});
db.createCollection("lottolog",{size:629145600, capped:true, max:3000000,autoIndexId:false});
db.createCollection("maintenancelog",{size:10485760, capped:true, max:3000,autoIndexId:false});
