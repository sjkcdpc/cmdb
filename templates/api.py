import urllib2
try:
    import json
except:
    import simplejson as json

url ="http://127.0.0.1:8130/api/v1/hosts/all"
data = urllib2.urlopen(url).read()
json_data = json.loads(data)
# print json_data
# print type(json_data)
for i in json_data.split('|'):
    print i
    print type(i)
    a = json.loads(i)
    print a['Os']

