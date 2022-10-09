#!/usr/bin/env bash

:<<!

title=check remote interface response
cid=0
pid=0

1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Validate its format >> `^[0-9]{8}`

!

resp=$(curl -s 'https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1')  # apt-get install curl if needed
elem=`echo $resp | grep -o '"startdate":"[^"]*"' | sed 's/^.*:"//g' | sed 's/"//g'`

echo "$elem"
