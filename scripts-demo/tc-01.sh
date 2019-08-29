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

echo '中文测试'

tm=`ping -c 1 zentao.com 2>/dev/null | grep 'time='`

echo '#'
echo $tm
