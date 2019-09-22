#!/usr/bin/env php

<?php
/**
[case]

title=expect with format string
cid=1
pid=1

step1 >> %s%d

[esac]
*/

checkStep1() || print(">> abc123\n");

function checkStep1(){}