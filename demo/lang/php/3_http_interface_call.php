#!/usr/bin/env php
<?php
/**

title=check remote interface response
cid=0
pid=0

1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Check its format >> `^[a-z0-9]{8}`

*/

$resp = file_get_contents('https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1');
$json = json_decode($resp);
echo $json->images[0]->startdate . "\n";
