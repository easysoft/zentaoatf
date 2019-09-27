#!/usr/bin/env php
<?php
/**
[case]
title=step multi lines
cid=0
pid=0

[1. group title 1]
  [1.1. steps]
    step 1.1
    step 1.2
  [1.1. expects]
    
  [1.2. steps]
    step 2.1
    step 2.2
  [1.2. expects]
    expect 2.1
    expect 2.2

[esac]
*/

if (checkStep1_2() || true) {
    print(">>\n");
    print("expect 2.1\n");
    print("expect 2.2\n");
}

function checkStep1_2(){}