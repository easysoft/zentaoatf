#!/usr/bin/env lua
--[[
title=check remote interface response
cid=0
pid=0

  1. Send a request to interface http://xxx
  2. Retrieve sessionID field from response json
  3. Check its format >> `^[a-z0-9]{26}`

]]

local http = require("socket.http") -- need luasocket library (luarocks install luasocket)
local ltn12 = require("ltn12")

function http.get(u)
   local t = {}
   local r, c, h = http.request{
      url = u,
      sink = ltn12.sink.table(t)
   }
   return r, c, h, table.concat(t)
end

r, c, h,body = http.get("http://max.demo.zentao.net/pms/?mode=getconfig")
if c~= 200 then
    print("ERR: " .. c)
else
    _, _, src = string.find(body, '"sessionID":"(.-)"')
    print(src)
end
