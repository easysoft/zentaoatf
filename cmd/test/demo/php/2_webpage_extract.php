#!/usr/bin/env php
<?php

/**

title=extract content from webpage
timeout=0
cid=2
pid=1

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `.*zt-logo.png`
1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `.*zt-logo.png`
1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `.*zt-logo.png`

*/

$resp = file_get_contents('http://max.demo.zentao.net/user-login-Lw==.html');
preg_match_all("/<img src=\"(.*)\" .*>/U", $resp, $matches);
echo $matches[1][0] . "\n";
print("abc123\n");