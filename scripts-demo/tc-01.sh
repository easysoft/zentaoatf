#!/usr/bin/env bash
:<<!
<<<TC

caseId:         1
productId:      0
title:          Test site response time
steps:          steps that begin with @ are checkpoints
   step1           type "ping zentao.com"
   @step2          check the output contains "ttl"

expects:
# @step2
.*ttl.*

TC;
!

:: print the line with ttl
tm=`ping -c 1 zentao.com 2>/dev/null | grep 'ttl'`

echo '#'
echo $tm
