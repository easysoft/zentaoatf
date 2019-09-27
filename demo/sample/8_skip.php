#!/usr/bin/env php
<?php
/**
[case]

title=skip
cid=-1
pid=0

[group]
  step 1 >> expect 1

[group title 3]
  step 3.1 >> expect 1.1
  step 3.2 >> expect 1.2

[esac]
*/

checkPreCondition() || print("skip\n");
print(">> ...\n");

function checkPreCondition(){}