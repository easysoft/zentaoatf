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