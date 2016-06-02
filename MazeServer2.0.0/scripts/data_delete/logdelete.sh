#!/bin/sh
if [ $# == 0 ];then
echo "please input dir"
echo "Usage: $0 <dir>"
exit 1
fi
for dir in  $@
do
find $dir -mtime +7  -exec rm -f {} \;
done