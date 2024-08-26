#!/usr/bin/env php
<?php
/**

title=with multi groups
cid=1

- step1 @ expect 1
- step2 @ expect 2
- step3 @
  - step3.1 @ expect 3.1
  - step3.2 @ expect 3.2
*/

print("expect 1\n");
print("expect 2\n");
print("pass\n");

print("expect 3.1\n");
print("expect 3.2\n");
print("pass\n");
