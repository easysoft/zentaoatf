#!/usr/bin/env bash

:<<!
[case]

title=json api call
cid=1
pid=1

Send a request to interface http://xxx
Retrieve sessionID field from response json
Validate its format >> ^[a-z0-9]{26}

[esac]
!

resp=$(curl -s 'http://ruiyinxin.test.zentao.net/?mode=getconfig')
elem=`echo $resp | grep -o '"sessionID":"[^"]*"' | sed 's/^.*:"//g' | sed 's/"//g'`

echo $elem
