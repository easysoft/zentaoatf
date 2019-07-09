#!/usr/bin/env php 
<?php
include 'init.php';
title('a failed case example.');

run(printHello(chr(rand(ord('a'), ord('z'))))) && expect('hello world %i');
run(printHello(chr(rand(ord('a'), ord('z'))))) && expect('hello world [a-z]{2}');
run(printHello(time())) && expect('hello world [0-9]{11}');

function printHello($dynamic)
{
    return 'hello world ' . $dynamic;
}
