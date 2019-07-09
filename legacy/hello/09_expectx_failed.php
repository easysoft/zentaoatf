#!/usr/bin/env php 
<?php
<<<TC
title: a failed case with regular rules.
expectx:| 
    hello world [0-9]*.
    hello world [A-Z]*.
TC;

echo 'hello world ' . time() . ".\n" ;
echo 'hello world ' . chr(rand(ord('a'), ord('z'))) . ".\n";
