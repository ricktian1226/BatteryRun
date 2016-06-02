#!/bin/sh
if [ $# != 1 ];then
    echo "Usage : $0 <uid>"
	echo "e.g: $0 1427450118065046536847"
	exit 1
fi

uid=$1

./usercommondelete.sh m.br2:27017 $uid
./userdelete.sh m.br2:27017 $uid
./userlogdelete.sh m.br2:27017 $uid
