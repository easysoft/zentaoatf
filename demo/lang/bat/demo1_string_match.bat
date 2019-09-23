#!/usr/bin/env bash

:<<!
[case]

title=check string matches pattern
cid=0
pid=0

exactly match            >> abc123
regular expression match >> abc\d{3}
format string match      >> %s%d

[esac]
!

echo ">> hello"
echo ">> 13905120512"
echo ">> abc123"
