#!/usr/bin/env tclsh
set case {
[case]
title=extract content from webpage
cid=0
pid=0

[group]
  1. Load web page from url http://xxx 
  2. Retrieve img element zt-logo.png in html 
  3. Check img exist >> `.*zt-logo.png`

[esac]
}

package require http

set url http://pms.zentao.net/user-login.html

set http [::http::geturl $url]
set html [::http::data $http]

regexp -- {<img src='(.*?)' .*>} $html match elem
puts ">> $elem"