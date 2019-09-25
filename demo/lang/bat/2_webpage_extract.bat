goto start
[case]

title=extract content from webpage
cid=0
pid=0

[group]
1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> .*zt-logo.png

[esac]
:start

@echo off

for /f "delims=" %%s in ('curl -s  "http://pms.zentao.net/user-login.html" ^| findstr/irc:"<img src='[^']*'.*/>"') do (
    for /f "tokens=2 delims='" %%i in ("%s%") do (
        set b=%%~i
    )
    echo %b%
)
