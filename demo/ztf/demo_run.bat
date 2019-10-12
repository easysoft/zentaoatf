@echo off

rmdir /s/q log
ztf.exe run demo/lang demo/sample
ztf.exe run demo/ztf/demo_check.bat