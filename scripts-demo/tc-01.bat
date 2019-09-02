goto start
<<<TC

caseId:         1
productId:      0
title:          Test network connection
steps:          steps that begin with @ are checkpoints
   step1           type "ping zentao.com"
   @step2          check the output contains "TTL"

expects:
# @step2
.*TTL.*

TC;
:start

@echo off

:: print the line with TTL
for /f "delims=" %%i in ('ping zentao.com -n 1 ^| findstr "TTL"') do set output=%%i

echo #
echo !output!
