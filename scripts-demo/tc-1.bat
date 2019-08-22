goto start

caseId:         1
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

:start

for /f "tokens=7" %%i in ('ping zentao.com ^| findstr /v "Active Proto"') do tm = %%i