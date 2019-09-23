#!/usr/bin/env php

<?php
/**
[case]

title=extract content from webpage
cid=0
pid=0

Load web page from url http://xxx
Find img element zt-logo.png in html >> .*zt-logo.png

[esac]
*/

$resp = file_get_contents('http://pms.zentao.net/user-login.html');
preg_match_all("/<img src='(.*)' .*>/U", $resp, $matches);
echo ">> " . $matches[1][0] . "\n";
