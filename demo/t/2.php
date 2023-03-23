#!/usr/bin/env php
<?php

/**

title=测试用例标题
cid=2

 - 步骤1 @
 - 步骤2
   - 子步骤2.1 @
   - 子步骤2.2 @
 - 步骤3
 - 步骤4 @

*/

print("@期待结果1\n");
print("@期待结果2.1\n");

print("@{\n");
print("期待结果2.2-1\n");
print("期待结果2.2-2\n");
print("}\n");

print("@期待结果3\n");

?>