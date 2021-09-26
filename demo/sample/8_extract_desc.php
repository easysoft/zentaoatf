#!/usr/bin/env php
<?php

/**

title=sync step from comments
cid=1
pid=0

1 >> expect 1

group2
  2.1 >> expect 2.1
  2.2 >> expect 2.2
  2.3 >> expect 2.3  

multi line expect >>
  expect 3.1
  expect 3.2
>>

4 >> expect 4
5 >> expect 5
step 6 >> expect 6
step 7 >> expect 7
step 8 >> expect 8
step 9 >> expect 9
step 10 >> expect 10

*/

print("expect 1\n"); // Step: 1  >> expect 1

/* group: group2 */
print("expect 2.1\n"); // step: 2.1  >> expect 2.1
print("expect 2.2\n"); // step: 2.2  >> expect 2.2
print("expect 2.3\n"); // step: 2.3  >> expect 2.3  ]]

/*
step: multi line expect >>
  expect 3.1
  expect 3.2
>>
*/
print(">>\n");
print("expect 3.1\n");
print("expect 3.2\n");
print(">>\n");

print("expect 4\n"); // step: 4  >> expect 4
print("expect 5\n"); // step: 5  >> expect 5

p("expect 6\n") && e('expect 6'); // step 6
p("expect 7\n") && e('expect 7'); # step 7
// step 8
p("expect 8\n") && e('expect 8');
# step 9
p("expect 9\n") && e('expect 9');
/** step 10 */
p("expect 10\n") && e('expect 10');

function p($msg) {
    print($msg);
    return true;
}
function e($msg) {
}