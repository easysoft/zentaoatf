#!/usr/bin/env php

<?php
/**
[case]

title=check remote interface response
cid=0
pid=0

Send a request to interface http://xxx
Retrieve sessionID field from response json
Validate its format >> ^[a-z0-9]{26}

[esac]
*/

$resp = file_get_contents('http://pms.zentao.net?mode=getconfig');
$json = json_decode($resp);
echo ">> " . $json->sessionID . "\n";
