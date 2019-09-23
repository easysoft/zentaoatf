#!/usr/bin/env tclsh
set case {
[case]

title=check string matches pattern
cid=0
pid=0

[group]
1. exactly match            >> hello
2. regular expression match >> 1\d{10}
3. format string match      >> %s%d

[esac]
}

puts ">> hello"
puts ">> 13905120512"
puts ">> abc123"
