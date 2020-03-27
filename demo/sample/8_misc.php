#!/usr/bin/env php
<?php
/**
[case]
title=step multi lines
cid=0
pid=0

[group]
  1. step 1 >> expect 1
  2. step 2 >> expect 2

[3. group title 1]
  [3.1. steps]
    step 3.1.1
    step 3.1.2
  [3.1. expects]
    
  [3.2. steps]
    step 3.2.1
    step 3.2.2
  [3.2. expects]
    expect 3.2.1
    expect 3.2.2

[group]
  4. step 4
  5. step 5 >> expect 5

[esac]
*/

if (checkStep1_2() || true) {
    print(">> expect 1\n");
    print(">> expect 2\n");

    print(">>\n");
    print("expect 3.2.1\n");
    print("expect 3.2.2\n");

    print(">> expect 5\n");
}

function checkStep1_2(){}