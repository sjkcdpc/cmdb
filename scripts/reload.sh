#!/usr/bin/env bash

echo "Current running:"
ps aux | grep -v grep | grep cmdb
count=`ps aux | grep -v grep | grep cmdb | wc -l`
if [ ${count} -gt 0 ]; then
    echo "reload all running apps gracefully:"
    ps aux | grep -v grep | grep cmdb | awk '{print $2}' | xargs kill -SIGUSR2
else
    echo "start new one:"
    nohup ./../cmdb start &
fi
