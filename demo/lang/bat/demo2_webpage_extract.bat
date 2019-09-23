#!/usr/bin/env bash

:<<!
[case]

title=extract content from webpage
cid=0
pid=0

Load web page from url http://xxx
Find img element zt-logo.png in html >> .*zt-logo.png

[esac]
!

resp=$(curl -s http://pms.zentao.net/user-login.html)
elem=`echo $resp | grep -o '<img[^>]*src="[^"]*"' | grep -o '[^"]*.png'`

echo $elem
