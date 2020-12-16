#!/usr/bin/env php
<?php
/**
[case]
title=the simple demo for ztf
cid=0
pid=0

[group]
  1. step1 >> expect好 1
  2. step2
  3. step3 >> expect 3

[esac]
*/

checkStep1() || print(">> - expect好 1\n");
checkStep3() || print(">> expect 3\n");

function checkStep1(){}
function checkStep3(){}