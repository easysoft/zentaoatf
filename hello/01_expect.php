#!/usr/bin/env php 
<?php
include 'init.php';
title('the hello word case for php scripts');
run(sprintf('helloworld')) && expect('helloworld');
