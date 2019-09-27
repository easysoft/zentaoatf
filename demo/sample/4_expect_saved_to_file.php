#!/usr/bin/env php
<?php
/**
[case]
title=step multi_lines
cid=0
pid=0

[group]
  1. step 1 
  2. step 2 

[3. group title 3]
  [3.1. steps]
    step 3.1
    step 3.2
  [3.1. expects]
    >>

[esac]
*/

checkStep1() || print(">> expect 1\n");

if (checkStep3() || true) {
    print(">>\n");
    print("expect 3.1\n");
    print("expect 3.2\n");
}

function checkStep1(){}
function checkStep3(){}