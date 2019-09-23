#!/usr/bin/env php

<?php
/**
[case]

title=the simple demo for ztf
cid=1
pid=1

Send a request to interface http://xxx
Retrieve sessionID field from response json
Validate its format >> ^[a-z0-9]{26}

[esac]
*/

$resp = file_get_contents('http://ruiyinxin.test.zentao.net/?mode=getconfig');
$json = json_decode($resp);
echo ">> " . $json->sessionID . "\n";
