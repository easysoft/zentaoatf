#!/usr/bin/env php
<?php
/**
[case]
title=expect with regx by using backquote `
cid=0
pid=0

[group]
  1. mobile phone >> `^1[0-9]\d{9}$`
  2. email        >> `^.+@.+\..+$`
  3. web url      >> `^https?.+\..+`
  4. ip address   >> `(\d+)\.(\d+)\.(\d+)\.(\d+)`

[esac]
*/

print(">> 13912345678\n");
print(">> 462826@qq.com\n");
print(">> https://www.zentao.net/index.html\n");
print(">> 192.168.0.1/24\n");
