#!/bin/bash

# https://phoenixnap.com/kb/set-up-cron-job-linux
# https://unix.stackexchange.com/questions/257960/how-do-i-find-files-older-than-1-days-using-mtime/257966
# https://ostechnix.com/how-to-find-and-delete-files-older-than-x-days-in-linux/
#find . -daystart -mtime +0 -print
find /root/*.xlsx -mtime +7 -delete
find /root/trans/*.xlsx -mtime +7 -delete
find /root/asin/*.jpg -mtime +7 -delete