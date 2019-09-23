#!/usr/bin/env bash

:<<!
[case]

title=webpage extract
cid=1
pid=1

Load web page from url http://xxx
Find img element zt-logo.png in html >> .*zt-logo.png

[esac]
!

resp=$(curl -s http://ruiyinxin.test.zentao.net/user-login.html)
elem=`echo $resp | grep -o '<img[^>]*src="[^"]*"' | grep -o '[^"]*.png'`

echo $elem
