goto start
[case]

title=check string matches pattern
cid=0
pid=0

[group]
1. exactly match            >> hello
2. regular expression match >> 1d{10}
3. format string match      >> %s%d

[esac]
:start

@echo off

echo ">> hello"
echo ">> 13905120512"
echo ">> abc123"
