#!/usr/bin/env bash

:<<!
[case]
title=test demo scripts
cid=-1
pid=0

[group]
  1. javascript/1_string_match.js >> Pass
  2. javascript/2_webpage_extract.js >> Pass
  3. javascript/3_http_inteface_call.js >> Pass
  4. lua/1_string_match.lua >> Pass
  5. lua/2_webpage_extract.lua >> Pass
  6. lua/3_http_inteface_call.lua >> Pass
  7. perl/1_string_match.pl >> Pass
  8. perl/2_webpage_extract.pl >> Pass
  9. perl/3_http_inteface_call.pl >> Pass
  10. php/1_string_match.php >> Pass
  11. php/2_webpage_extract.php >> Pass
  12. php/3_http_inteface_call.php >> Pass
  13. python/1_string_match.py >> Pass
  14. python/2_webpage_extract.py >> Pass
  15. python/3_http_inteface_call.py >> Pass
  16. ruby/1_string_match.rb >> Pass
  17. ruby/2_webpage_extract.rb >> Pass
  18. ruby/3_http_inteface_call.rb >> Pass
  19. shell/1_string_match.sh >> Pass
  20. shell/2_webpage_extract.sh >> Pass
  21. shell/3_http_inteface_call.sh >> Pass
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
!

cat log/001/result.txt | grep '^(' | while read line
do
   echo '>>' $line
done