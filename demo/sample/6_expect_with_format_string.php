#!/usr/bin/env php
<?php
/**

title=expect with format string by using backquote `
cid=0
pid=0

string              >> `^%s$`
integer             >> `^%d$`
signed integer      >> `^%i$`
hexadecimal number  >> `^%x$`
float               >> `^%f$`
char                >> `^%c$`

*/

print("abc\n");
print("123\n");
print("-123\n");
print("0XAF\n");
print("+1.23\n");
print("x\n");
