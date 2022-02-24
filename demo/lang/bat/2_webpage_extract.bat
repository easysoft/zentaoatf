@echo off
setlocal enabledelayedexpansion
goto start

title=extract content from webpage
cid=0
pid=0

  1. Load web page from url http://xxx
  2. Retrieve img element zt-logo.png in html
  3. Check img exist >> `.*.png`

:start

for /f "delims=" %%a in ('curl -s  "http://max.demo.zentao.net/user-login-Lw==.html" ^| findstr/irc:"<img src="') do (
    set var=%%a
)
for /f "tokens=1 delims='" %%i in ("!var!") do (
    set var2=%%~i
)
echo !var2!
