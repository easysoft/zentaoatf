#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''
[case]

title=extract content from webpage
cid=0
pid=0

[group]
1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `.*zt-logo.png`

[esac]
'''

import requests
import re

html = requests.get('http://pms.zentao.net/user-login.html').content #need requests library (pip/pip3 install requests)
elem = re.search(r"<img src='(.*?)' .*>", html.decode("utf-8")).group(1)
print('>> ' + elem)
