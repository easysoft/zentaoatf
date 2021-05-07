#!/usr/bin/env php
<?php
/**

title=simple demo
cid=1
pid=1

1. step1 >> expect 1
2. step2
3. step3 >> expect 3

*/

checkStep1() || print("expect 1\n");
checkStep3() || print("expect 3\n");

function checkStep1(){}
function checkStep3(){}
