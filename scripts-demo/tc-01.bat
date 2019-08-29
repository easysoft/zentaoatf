goto start
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
:start

@echo off
Setlocal enabledelayedexpansion

for /f "delims=" %%i in ('ping zentao.com -n 1 ^| findstr "TTL"') do set output=%%i

echo #
echo !output!
