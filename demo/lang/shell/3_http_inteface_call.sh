#!/usr/bin/env bash

:<<!
[case]

title=check remote interface response
cid=0
pid=0

[group]
1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Validate its format >> ^[a-z0-9]{26}

[esac]
!

resp=$(curl -s 'http://pms.zentao.net?mode=getconfig')
elem=`echo $resp | grep -o '"sessionID":"[^"]*"' | sed 's/^.*:"//g' | sed 's/"//g'`

echo ">> $elem"
