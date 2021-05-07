#!/usr/bin/env php
<?php
/**

[case]
title=with multi groups
cid=0
pid=0

[group]
  1. step 1 >> expect 1
  2. step 2

[3. group title 3]
  3.1 step
  3.2 step >> expect 3

[4. group title 1]
  [4.1. steps]
    step 4.1.1
    step 4.1.2
  [4.1. expects]

  [4.2. steps]
    step 4.2.1
    step 4.2.2
  [4.2. expects]
    expect 4.2.1
    expect 4.2.2

[esac]

*/

print(">>expect 1\n");
print(">>expect 3\n");

print(">>\n");
print("expect 4.2.1\n");
print("expect 4.2.2\n");
