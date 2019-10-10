@echo off
goto start
[case]
title=check string matches pattern
cid=0
pid=0

[group]
  1. bat/1_string_match.sh >> Pass
  2. bat/2_webpage_extract.sh >> Pass
  3. bat/3_http_inteface_call.sh >> Pass
  4. javascript/1_string_match.js >> Pass
  5. javascript/2_webpage_extract.js >> Pass
  6. javascript/3_http_inteface_call.js >> Pass
  7. lua/1_string_match.lua >> Pass
  8. lua/2_webpage_extract.lua >> Pass
  9. lua/3_http_inteface_call.lua >> Pass
  10. perl/1_string_match.pl >> Pass
  11. perl/2_webpage_extract.pl >> Pass
  12. perl/3_http_inteface_call.pl >> Pass
  13. php/1_string_match.php >> Pass
  14. php/2_webpage_extract.php >> Pass
  15. php/3_http_inteface_call.php >> Pass
  16. python/1_string_match.py >> Pass
  17. python/2_webpage_extract.py >> Pass
  18. python/3_http_inteface_call.py >> Pass
  19. ruby/1_string_match.rb >> Pass
  20. ruby/2_webpage_extract.rb >> Pass
  21. ruby/3_http_inteface_call.rb >> Pass
  22. tcl/1_string_match.tl >> Pass
  23. tcl/2_webpage_extract.tl >> Pass
  24. tcl/3_http_inteface_call.tl >> Pass
  25. sample/1_simple.php >> Pass
  26. sample/2_with_group.php >> Pass
  27. sample/3_step_multi_lines.php >> Pass
  28. sample/4_expect_saved_to_file.php >> Pass
  29. sample/5_expect_with_express.php >> Pass
  30. sample/6_expect_with_regx.php >> Pass
  31. sample/7_expect_with_format_string.php >> Pass
  32. sample/8_misc.php >> Pass
  33. sample/9_skip.php >> Skip

[esac]
:start

for /f "delims="  eol^=^( %%i in (log/001/result.txt) do (
    set "a=%%i"
    echo ^>^> !a!
)
