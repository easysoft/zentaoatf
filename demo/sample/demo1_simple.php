#!/usr/bin/env php

<?php
/**
[case]

title=the simple demo for ztf
cid=1
pid=1

step1 >> expect 1
step2
step3 >> expect 3

[esac]
*/

checkStep1() || print(">> expect 1\n");
checkStep3() || print(">> expect 3\n");

function checkStep1(){}
function checkStep3(){}
