#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
[case]

title=check remote interface response
cid=0
pid=0

[group]
1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Check its format >> `^[a-z0-9]{26}`

[esac]
'''

import requests
import json

jsonStr = requests.get('http://pms.zentao.net/?mode=getconfig').content #need requests library (pip/pip3 install requests)
jsonObj = json.loads(jsonStr)
print('>> ' + jsonObj['sessionID'])
