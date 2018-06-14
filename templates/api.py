#!/usr/bin/env python
# coding=UTF-8


import urllib2

try:
    import json
except:
    import simplejson as json

url = "http://127.0.0.1:8130/api/v1/cmdb/hosts"
data = urllib2.urlopen(url).read()
json_data = json.loads(data)

for i in range(len(json_data)):
    line_info = json_data[i]
    json_info = json.JSONDecoder().decode(line_info)
    print(json_info['LanIp'])