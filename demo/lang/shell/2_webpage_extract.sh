#!/usr/bin/env bash

:<<!

title=extract content from webpage
cid=0
pid=0

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `.*zt-logo.png`

!

resp=$(curl -s http://max.demo.zentao.net/user-login-Lw==.html)   # apt-get install curl if needed
elem=`echo $resp | grep -o '<img[^>]*src="[^"]*"'`

echo "$elem"
