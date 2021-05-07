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

[1. group title 1]
  [1.1. steps]
    step 1.1.1
    step 1.1.2
  [1.1. expects]

  [1.2. steps]
    step 1.2.1
    step 1.2.2
  [1.2. expects]
    expect 1.2.1
    expect 1.2.2

[esac]

*/

print(">>expect 1\n");
print(">>expect 3\n");

print(">>\n");
print("expect 1.2.1\n");
print("expect 1.2.2\n");
