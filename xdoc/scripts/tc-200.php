<?php
<<<TC
title: 登录失败账号锁定策略
caseId: 200
steps: // @开头的为含验证点的步骤
   step2000 // 连续输入3次错误的密码
   @step2010 // 第4次尝试登录
   group2100 // 不连续输入3次错误的密码
      step2101 // 输入2次错误的密码
      step2102 // 输入1次正确的密码
      step2103 // 再输入1次错误的密码
      @step2104 // 再输入1次正确的密码

expects:
#
111
222
#
aaa
bbb

TC;

// 此处编写操作步骤代码

echo "#\n";   // 验证点@step2010的标记位，请勿删除
echo "111\n";
echo "222\n";

echo "#\n";   // 验证点@step2104的标记位，请勿删除
echo "aaa\n";
echo "bbb\n";


?>
