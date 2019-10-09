@echo off
goto start
[case]
title=check string matches pattern
cid=0
pid=0

[group]
javascript/1_string_match.js             >> Pass
javascript/2_webpage_extract.js          >> Pass
javascript/3_http_inteface_call.js       >> Pass
lua/1_string_match.lua                   >> Pass
lua/2_webpage_extract.lua                >> Pass
lua/3_http_inteface_call.lua             >> Pass
perl/1_string_match.pl                   >> Pass
perl/2_webpage_extract.pl                >> Pass
perl/3_http_inteface_call.pl             >> Pass
php/1_string_match.php                   >> Pass
php/2_webpage_extract.php                >> Pass
php/3_http_inteface_call.php             >> Pass
python/1_string_match.py                 >> Pass
python/2_webpage_extract.py              >> Pass
python/3_http_inteface_call.py           >> Pass
ruby/1_string_match.rb                   >> Pass
ruby/2_webpage_extract.rb                >> Pass
ruby/3_http_inteface_call.rb             >> Pass
shell/1_string_match.sh                  >> Pass
shell/2_webpage_extract.sh               >> Pass
shell/3_http_inteface_call.sh            >> Pass
tcl/1_string_match.tl                    >> Pass
tcl/2_webpage_extract.tl                 >> Pass
tcl/3_http_inteface_call.tl              >> Pass
sample/1_simple.php                      >> Pass
sample/2_with_group.php                  >> Pass
sample/3_step_multi_lines.php            >> Pass
sample/4_expect_saved_to_file.php        >> Pass
sample/5_expect_with_express.php         >> Pass
sample/6_expect_with_regx.php            >> Pass
sample/7_expect_with_format_string.php   >> Pass
sample/8_misc.php                        >> Pass
sample/9_skip.php                        >> Pass

[esac]
:start

REM 此处编写代码，读取log，输出内容到控制台，ztf框架回同上述步骤进行验证