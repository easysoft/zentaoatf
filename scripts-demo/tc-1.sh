#!/usr/bin/env bash
<<TC
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

TC

timeout=300

((count = 3)) #Number to test

while [[ $count -ne 0 ]] ; do
    #get time field
    tm=`ping -c 1 zentao.com 2>/dev/null | grep 'time=' | sed 's/.*time=\([.0-9]*\) ms/\1/g' | awk -F. '{print $1}'`
    echo $tm

    if [[ $tm -gt $timeout ]] ; then #timeout
        ((count = 1)) # break
    fi
    ((count = count - 1))
done

echo '#' #checkpoint 1

if [ ! -n "$tm" ]; then
    echo 'unknown'
elif [[ $tm -gt $timeout ]]; then
    echo 'timeout'
else
    echo 'work'
fi