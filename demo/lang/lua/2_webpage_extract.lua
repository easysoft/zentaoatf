#!/usr/bin/env lua
--[[
title=extract content from webpage
cid=0
pid=0

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `必应`

]]

require("socket")
local https = require("ssl.https")  -- need luasocket, luasec modules
local body, code, headers, status = https.request("https://cn.bing.com")

_, _, src = string.find(body, "<title>(.-)</title>")
print(src)
