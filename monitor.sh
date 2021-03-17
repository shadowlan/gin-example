#!/bin/bash
# 监控main进程是否正常运行，如果没有则重新运行
ps aux | grep "./main" | grep -v grep
if [ $? != 0 ];then
   d=$(date +"%m%d%H%M")
   mv mainout.txt  mainout.txt."$d"
   nohup ./main >mainout.txt 2>&1 &
fi