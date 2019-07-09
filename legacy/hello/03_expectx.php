#!/usr/bin/env php 
<?php
include 'init.php';
title('with regular rules');
run(sprintf('hello world ' . time())) && expectx('hello world [0-9]*');
