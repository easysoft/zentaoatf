#!/usr/bin/env python
'''
[case]

title=extract content from webpage
cid=0
pid=0

[group]
1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> .*zt-logo.png

[esac]
'''

import requests
import re

html = requests.get('http://pms.zentao.net/user-login.html').content #need to install luasocket library, easy_install requests
elem = re.search(r"<img src='(.*?)' .*>", html).group(1)
print('>> ' + elem)

