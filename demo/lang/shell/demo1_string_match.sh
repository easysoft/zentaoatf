#!/usr/bin/env bash

<?php
/**
[case]

title=string match
cid=1
pid=1

test string             >> abc123
test regular expression >> abc\d{3}
test format string      >> %s%d

[esac]
*/

// your business logic here
$str = "abc" . "123";

print(">> $str\n");
print(">> $str\n");
print(">> $str\n");
