#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''

title=extract content from webpage
cid=0
pid=0

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `.*zt-logo.png`

'''

import requests
import re

html = requests.get('http://max.demo.zentao.net/user-login-Lw==.html').content #need requests library (pip/pip3 install requests)
elem = re.search(r"<img src=\"(.*?)\" .*>", html.decode("utf-8")).group(1)
print(elem)
