#!/usr/bin/env lua
--[[
title=extract content from webpage
cid=0
pid=0

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `.*.png`

]]

local http = require("socket.http") -- need luasocket library (luarocks install luasocket)
local ltn12 = require("ltn12")

function http.get(u)
   local t = {}
   local r, c, h = http.request{
      url = u,
      sink = ltn12.sink.table(t)}
   return r, c, h, table.concat(t)
end

r,c,h,body = http.get("http://max.demo.zentao.net/user-login-Lw==.html")
if c~= 200 then
    print("ERR: " .. c)
else
    _, _, src = string.find(body, "<img%ssrc=\"(.-)\" .*>")
    print(src)
end
