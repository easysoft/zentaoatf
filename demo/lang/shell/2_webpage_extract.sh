#!/usr/bin/env bash

:<<!
[case]

title=extract content from webpage
cid=0
pid=0

[group]
1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> .*zt-logo.png

[esac]
!

resp=$(curl -s http://pms.zentao.net/user-login.html)   # apt-get install curl if needed
elem=`echo $resp | grep -o "<img[^>]*src='[^']*'" | grep -o "[^']*.png"`

echo ">> $elem"
