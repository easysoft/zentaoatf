#!/usr/bin/env php
<?php
/**

title=step multi lines
cid=0
pid=0

steps
  step 1.1
  step 1.11 >>
  step 1.2 >>
    expect 1.2 line 1
    expect 1.2 line 2
  >>

*/
print("pass\n");

print(">>\n");
print("expect 1.2 line 1\n");
print("expect 1.2 line 2\n");
print(">>\n");
