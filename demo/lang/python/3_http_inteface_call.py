#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''

title=check remote interface response
cid=0
pid=0

1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Check its format >> `^[a-z0-9]{26}`

'''

import requests
import json

jsonStr = requests.get('http://max.demo.zentao.net/pms/?mode=getconfig').content #need requests library (pip/pip3 install requests)
jsonObj = json.loads(jsonStr)
print(jsonObj['sessionID'])
