#!/bin/bash
configDir=/home/xyao/workspace/MazeServer2.0.0/conf
cd $configDir
svn up
cd /home/xyao/workspace/MazeServer2.0.0/gocode
./build_all.sh $1
cp $configDir/battery_*_server.ini $configDir/gnats*.conf ./
cp certs2/aps_superbman_*.pem ./
#tar zcvf fullPackage.tar.gz ../conf/*.ini ../conf/gnats*.conf $1
tar zcvf fullPackage.tar.gz battery_*_server.ini gnats*.conf gnatsd $1 httpdoc aps_superbman_*.pem


