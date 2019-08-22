#!/usr/bin/env bash
<<TC
title:shell hello world.
expect:hello world.
TC

timeout=300

((count = 3)) #  Number to test

while [[ $count -ne 0 ]] ; do
    // get time field
    tm=`ping -c 1 zentao.com 2>/dev/null | grep 'time=' | sed 's/.*time=\([.0-9]*\) ms/\1/g' | awk -F. '{print $1}'`
    echo $tm

    if [[ $tm -gt $timeout ]] ; then // timeout
        ((count = 1)) // break
    fi
    ((count = count - 1))
done

if [ ! -n "$tm" ]; then
    echo 'Unknown'
elif [[ $tm -gt $timeout ]]; then
    echo 'Timeout'
else
    echo 'Work'
fi