#!/usr/bin/env python
'''
[case]

title=check string matches pattern
cid=1
pid=1

[group]
1. exactly match            >> hello
2. regular expression match >> 1\d{10}
3. format string match      >> %s%d

[esac]
'''

print(">> hello")
print(">> 13905120512")
print(">> abc123")
