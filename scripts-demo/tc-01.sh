#!/usr/bin/env bash

:<<!
<<<TC

caseId:         1
caseIdInTask:   0
taskId:         0
title:          Test site response time
steps:          steps that begin with @ are checkpoints
   step1           type "ping zentao.com" to send ICMP request
   step2           get response time from output
   step3           if time > 300ms, break the cycle
                      time < 300ms, continue
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
!

timeout=500

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

echo '#' #checkpoint start

if [ ! -n "$tm" ]; then
    echo 'unknown'
elif [[ $tm -gt $timeout ]]; then
    echo 'timeout'
else
    echo 'pass'
fi