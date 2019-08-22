#!/usr/bin/env bash
<<TC
caseId:         1
caseIdInTask:   0
taskId:         0
title:          Demo测试用例
steps:          @开头的为含验证点的步骤
   step1           ping服务器zentao.com
   @step2           验证有返回，且响应时间小于300ms

expects:
# @step2
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