#!/usr/bin/env php
<?php

/**

title=sync step from comments
cid=1
pid=0

group1
1.1 >> expect 1.1
1.2 >> expect 1.2
1.3 >> expect 1.3  ]]

group2
2.1 >> expect 2.1
2.2 >> expect 2.2
2.3 >> expect 2.3  ]]

multi line expect >>
expect 3.1
expect 3.2
>>

4 >> expect 4
5 >> expect 5

*/


/* group: group1 */
// Step: 1.1    >> expect 1.1
// Step: 1.2    >> expect 1.2
// Step: 1.3    >> expect 1.3  ]]

/* group: group2 */
// Step: 2.1    >> expect 2.1
// Step: 2.2    >> expect 2.2
// Step: 2.3    >> expect 2.3  ]]

/**
step: multi line expect >>
expect 3.1
expect 3.2
>>
 */

// step: 4 >> expect 4
// step: 5 >> expect 5

print("expect 1.1\n");
print("expect 1.2\n");
print("expect 1.3\n");

print("expect 2.1\n");
print("expect 2.2\n");
print("expect 2.3\n");

print(">>\n");
print("expect 3.1\n");
print("expect 3.2\n");
print("<<\n");

print("expect 4\n");
print("expect 5\n");