#!/usr/bin/env bash

echo "Current running:"

ps aux | grep -v grep | grep cmdb

killall cmdb
