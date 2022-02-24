#!/usr/bin/env php
<?php
/**

title=simple demo
cid=1
pid=1

step1 >> expect 1
step2 >>
step3 >> expect 3

*/

checkStep1() || print("expect 1\n");
print("pass\n");
checkStep3() || print("expect 3\n");

function checkStep1(){}
function checkStep3(){}

stdErr('it is stderr msg');

function stdErr($msg) {
    fwrite(STDERR, "$msg\n");
}
