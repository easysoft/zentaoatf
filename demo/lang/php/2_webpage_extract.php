#!/usr/bin/env php
<?php
/**

title=extract content from webpage
cid=0
pid=0

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `必应`

*/

$resp = file_get_contents('https://cn.bing.com');
preg_match_all("/<title>(.*)</U", $resp, $matches);
echo $matches[1][0] . "\n";
