1、userdeleter.sh是用于删除某一指定玩家数据的脚本工具，使用方法如下：
./userdeleter.sh 192.168.93.129:27017 142256645185451
2、deletedynamicdata.sh是用于删除所有动态数据的脚本工具，使用方法如下：
./deletedynamicdata.sh 192.168.93.129:27017

【注意】
如果mongo程序不是安装在/usr/bin目录下，需要修改下脚本中mongoExe的赋值。