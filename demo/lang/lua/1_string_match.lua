#!/usr/bin/env lua
--[[
title=check string matches pattern
cid=0
pid=0

[group]
  1. exactly match >> hello
  2. regular expression match >> `1\d{10}`
  3. format string match >> `%s%d`

]]

print("hello");
print("13905120512");
print("abc123");
