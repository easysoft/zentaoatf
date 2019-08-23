@echo off
Setlocal enabledelayedexpansion

set timeout=500

for %%a in (1,2,3) do (
	for /f "tokens=5" %%i in ('ping zentao.com -n 1 ^| findstr "TTL"') do set tmstr=%%i
	echo !tmstr!

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
	echo work
)
