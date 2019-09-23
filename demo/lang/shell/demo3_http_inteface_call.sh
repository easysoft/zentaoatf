#!/usr/bin/env bash

:<<!
[case]

title=check remote interface response
cid=0
pid=0

Send a request to interface http://xxx
Retrieve sessionID field from response json
Validate its format >> ^[a-z0-9]{26}

[esac]
!

resp=$(curl -s 'http://pms.zentao.net?mode=getconfig')
elem=`echo $resp | grep -o '"sessionID":"[^"]*"' | sed 's/^.*:"//g' | sed 's/"//g'`

echo $elem
