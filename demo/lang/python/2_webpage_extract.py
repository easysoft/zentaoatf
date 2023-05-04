#!/usr/bin/env python3
# -*- coding: utf-8 -*-
'''

title=extract content from webpage
cid=0
pid=0

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `必应`

'''

import requests
import re

## for Chinese display
import sys,io,platform
if(platform.system()=='Windows'):
   import sys,io
   sys.stdout = io.TextIOWrapper(sys.stdout.buffer,encoding='utf8')

html = requests.get('https://cn.bing.com').content #need requests library (pip/pip3 install requests)
elem = re.search(r"<title>(.*?)<", html.decode("utf-8")).group(1)
print(elem)
