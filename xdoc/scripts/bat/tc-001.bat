goto start
<<<TC

caseId:         -1
caseIdInTask:   0
taskId:         0
title:          测试服务器响应时间
steps:          @开头的为含验证点的步骤
   step1           执行ping命令，向zentao.com服务器发送ICMP请求
   step2           获取返回消息，截取响应时间字段，转换成整数
   step3           判断响应时间，如果超过300ms，则放弃后续请求并返回
                               如果小于300ms，重复以上步骤累计3次
   @step4          验证最后返回的响应时间，如果小于300ms，打印"work"
                                        如果大于300ms，打印"timeout"
                                        如果为空，打印unknown表明为未到服务器

expects:
# @step4
work

readme:
- Logs of test scripts，must expects章节中#号标注的验证点需保持一致对应
- 脚本中CODE打头的注释需用代码替换
- 参考样例https://github.com/easysoft/zentaoatf/tree/master/xdoc/sample

TC;
:start

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
		goto ret
	)
)

:ret
echo #
if !tm2! GTR !timeout! (
	echo timeout
) else (
	echo work
)
