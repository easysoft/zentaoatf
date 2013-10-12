#!/usr/bin/env php 
<?php
include 'init.php';
title('with multi lines.');
run(multiline()) && expect("2\n3\n");

function multiline()
{
    $lines  = (1 + 1 . "\n");
    $lines .= (1 + 2 . "\n");
    return $lines;
}
