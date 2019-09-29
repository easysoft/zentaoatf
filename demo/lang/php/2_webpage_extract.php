#!/usr/bin/env php
<?php
/**
[case]
title=extract content from webpage
cid=0
pid=0

[group]
  1. Load web page from url http://xxx 
  2. Retrieve img element zt-logo.png in html 
  3. Check img exist >> .*zt-logo.png

[esac]
*/

$resp = file_get_contents('http://pms.zentao.net/user-login.html');
preg_match_all("/<img src='(.*)' .*>/U", $resp, $matches);
echo ">> " . $matches[1][0] . "\n";