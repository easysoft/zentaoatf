#!/usr/bin/env php 
<?php
include 'init.php';
title('with format chars(http://qa.php.net/phpt_details.php#expectf_section)');
run('hello world ' . time()) && expect('hello world %i');
run('hello world zentao')    && expect('hello world %s');
