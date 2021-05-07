#!/usr/bin/env php
<?php
/**

title=expect with format string by using backquote `
cid=0
pid=0

1. string              >> `^%s$`
2. integer             >> `^%d$`
3. signed integer      >> `^%i$`
4. hexadecimal number  >> `^%x$`
5. float               >> `^%f$`
6. char                >> `^%c$`

*/

print("abc\n");
print("123\n");
print("-123\n");
print("0XAF\n");
print("+1.23\n");
print("x\n");
