@echo off
setlocal enabledelayedexpansion
goto start

title=check remote interface response
cid=0
pid=0

1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Validate its format >> `^[a-z0-9]{26}`

:start

for /f "tokens=9 delims=," %%a in ('curl -s  "http://max.demo.zentao.net/pms/?mode=getconfig"') do (
    set var=%%a
)

for /f "tokens=2 delims=:" %%i in ("!var!") do (
    set s=%%i
    set var2=!s:~1,26!
)

echo !var2!
