#!/usr/bin/env php
<?php
/**
[case]

title=with multi groups
cid=0
pid=0

[group]
  step 1 >> expect 1
  step 2 >> expect 2

[group title 3]
  step 3.1 >> expect 3.1
  step 3.2 >> expect 3.2
  
[esac]
*/

checkStep1() || print(">> expect 1\n");
checkStep2() || print(">> expect 2\n");

checkStep3_1() || print(">> expect 3.1\n");
checkStep3_2() || print(">> expect 3.2\n");

function checkStep1(){}
function checkStep2(){}
function checkStep3_1(){}
function checkStep3_2(){}