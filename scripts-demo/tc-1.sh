goto s
<<<TC

caseId:         1
caseIdInTask:   0
taskId:         0
title:          Test site response time
steps:          steps that begin with "@" have checkpoints
   step1           type "ping zentao.com" to send ICMP request
   step2           get response time from command line output
   step3           if time > 300ms, break the cycle
                      time < 300ms, continue totally 3 times
   @step4          check the last response time，if time < 300ms，print "pass"
                                                    time > 300ms，print "timeout"

expects:
# @step4
pass

readme:
- Print '#' in test log to match up with the ones in expects section
- Write test scripts to replace the lines begin with 'CODE'
- More examples, pls refer to https://github.com/easysoft/zentaoatf/tree/master/xdoc/sample

TC;
:s

@echo off
Setlocal enabledelayedexpansion
::chcp 65001
::chcp 936

set timeout=500

for %%a in (1,2,3) do (
	for /f "tokens=5" %%i in ('ping zentao.com -n 1 ^| findstr "TTL"') do set tmstr=%%i
	REM echo !tmstr!

	for /f "tokens=2 delims='='" %%x in ('echo !tmstr!') do set tm=%%x
	set tm2=!tm:~0,-2!
	echo !tm2!

	if !tm2! GTR !timeout! (
		goto r
	)
)

:r
echo # ::checkpoint
if !tm2! GTR !timeout! (
	echo timeout
) else (
	echo pass
)
