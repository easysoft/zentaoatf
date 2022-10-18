#!/usr/bin/env lua
--[[
title=check remote interface response
cid=0
pid=0

  1. Send a request to interface http://xxx
  2. Retrieve sessionID field from response json
  3. Check its format >> `^[a-z0-9]{8}`

]]

require("socket")
local https = require("ssl.https") -- need luasocket, luasec modules
local body, code, headers, status = https.request("https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1")

_, _, src = string.find(body, '"startdate":"(.-)"')
print(src)
