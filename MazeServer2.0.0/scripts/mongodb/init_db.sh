#!/bin/bash
mongoExe=/usr/bin/mongo
#导入服务配置
$mongoExe m.br2:20143/brcommondb -u superUser -p superUser init_apiconfig.js
#创建公共表的索引
$mongoExe m.br2:20143/brcommondb -u superUser -p superUser init_index_common.js
#创建平台相关表的索引
$mongoExe s.br2:20144/briosdb -u superUser -p superUser init_index.js
$mongoExe s.br2:20144/brandroiddb -u superUser -p superUser init_index.js
#创建日志表
$mongoExe a.br2:20145/brlogdb init_logcollection.js
