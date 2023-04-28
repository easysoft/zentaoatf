#!/usr/bin/env python3
# -*- coding: utf-8 -*-
'''

title=check remote interface response
cid=0
pid=0

1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Check its format >> `^[a-z0-9]{8}`

'''

import requests
import json

## for Chinese display
import sys,io,platform
if(platform.system()=='Windows'):
   import sys,io
   sys.stdout = io.TextIOWrapper(sys.stdout.buffer,encoding='utf8')

jsonStr = requests.get('https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1').content
jsonObj = json.loads(jsonStr)
print(jsonObj['images'][0]['startdate'])
