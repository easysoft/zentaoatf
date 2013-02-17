#!/usr/bin/env php 
<?php
<<<TC
title: a failed case with format chars.
expect: hello world %i.
TC;

echo 'hello world ' . chr(rand(ord('a'), ord('z')));
