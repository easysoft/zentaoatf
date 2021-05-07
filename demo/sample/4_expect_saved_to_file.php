#!/usr/bin/env php
<?php
/**

title = expect in .exp file
cid=0
pid=0

1. step 1 >>
2. step 2

3. step 3
   >>

4. step 4
   >>
*/

print("expect 1\n");
print("expect 3\n");

// step 4: two expect lines in .exp file for single >> symbol in definition.
print(">>\n expect 4 line 1\n\n expect 4 line 2\n <<\n");
