#!/usr/bin/env php

<?php
/**
[case]

title=the simple demo for ztf
cid=1
pid=1

step1 >> .*zt-logo.png

[esac]
*/

$resp = file_get_contents('http://ruiyinxin.test.zentao.net/user-login.html');
preg_match_all('/<img src="(.*)" .*>/U', $resp, $matches);
echo ">> " . $matches[1][0] . "\n";
