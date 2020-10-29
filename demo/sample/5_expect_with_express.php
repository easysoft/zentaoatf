#!/usr/bin/env php
<?php
/**
[case]
title=the expect with regx
cid=0
pid=0

[group]
  1. step1 >> `abc\d{3}`

[esac]
*/

checkStep1() || print(">> abc123\n");

function checkStep1(){}