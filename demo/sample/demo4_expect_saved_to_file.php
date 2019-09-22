#!/usr/bin/env php

<?php
/**
[case]

title=step multi_lines
cid=1
pid=1

[group]
   step 1 >>
   step 2

[group title 3]
  [3. steps]
    step 3.1
    step 3.2
  [3. expects]
    >>

[esac]
*/

checkStep1() || print(">> expect 1\n");

if (checkStep3() || true) {
    print(">>\n");
    print("expect 3.1\n");
    print("expect 3.2\n");
}

function checkStep1(){}
function checkStep3(){}
