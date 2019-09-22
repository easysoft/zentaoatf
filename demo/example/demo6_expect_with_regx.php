#!/usr/bin/env php

<?php
/**
[case]

title=expect with regx
cid=1
pid=1

step1 >> ^abc\d{3}$

[esac]
*/

checkStep1() || print(">> abc123\n");

function checkStep1(){}
