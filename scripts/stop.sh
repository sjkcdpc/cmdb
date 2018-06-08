#!/usr/bin/env bash

echo "Current running:"

ps aux | grep -v grep | grep cmdb

count=`ps aux | grep -v grep | grep cmdb | wc -l`

#./../cmdb start

if [ ${count} -gt 0 ]; then

    echo "stop all running apps gracefully:"
    ps aux | grep -v grep | grep cmdb | awk '{print $2}' | xargs kill -TERM

    ps aux | grep -v grep | grep cmdb
fi
