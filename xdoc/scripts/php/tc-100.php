<?php
<<<TC
caseId:   100
title:    用户登录
steps:    @开头的为含验证点的步骤
   1000         打开登录页面
   1010         输入正确的用户名和密码
   @1020        点击'登录'按钮

expects:
# 
/* @step1020期望结果, 可以有多行 */

readme:
- 脚本输出日志和expects章节中，#号标注的验证点需保持一致对应
- 脚本中/* */标注的需用代码替换，//注解的为说明文字
- 参考样例https://github.com/easysoft/zentaoatf/tree/master/xdoc/sample

TC ;

/* 此处编写操作步骤代码 */

echo "#\n";  // @step1020: 用户成功登录
/* 输出验证点实际结果 */

?>
