#!/usr/bin/env bash

:<<!

title=check remote interface response
cid=0
pid=0

1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Validate its format >> `^[a-z0-9]{26}`

!

resp=$(curl -s 'http://max.demo.zentao.net/pms/?mode=getconfig')  # apt-get install curl if needed
elem=`echo $resp | grep -o '"sessionID":"[^"]*"' | sed 's/^.*:"//g' | sed 's/"//g'`

echo "$elem"
