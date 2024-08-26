#!/usr/bin/env bash

:<<!
title = simple demo
cid=1

step1 >> expect 1
step2 >>
step3 >> expect 3

!


echo "expect 1"
echo "pass"
echo "expect 3"
echo "it is stderr msg"
