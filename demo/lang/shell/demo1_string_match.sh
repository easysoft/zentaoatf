#!/usr/bin/env bash

:<<!
[case]

title=string match
cid=1
pid=1

test string             >> abc123
test regular expression >> abc\d{3}
test format string      >> %s%d

[esac]
!

str="abc123"

echo ">> ${str}"
echo ">> ${str}"
echo ">> ${str}"
