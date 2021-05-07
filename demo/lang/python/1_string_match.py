#!/usr/bin/env python
# -*- coding: utf-8 -*-
'''

title=check string matches pattern
cid=1
pid=1

1. exactly match            >> hello
2. regular expression match >> `1\d{10}`
3. format string match      >> `%s%d`
4. with Chinese      >> 中文

'''
## for Chinese display
import sys,io,platform
if(platform.system()=='Windows'):
   import sys,io
   sys.stdout = io.TextIOWrapper(sys.stdout.buffer,encoding='utf8')

print("hello")
print("13905120512")
print("abc123")

print("中文")
