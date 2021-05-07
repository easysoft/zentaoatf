#!/usr/bin/env tclsh
set case {

title=extract content from webpage
cid=0
pid=0

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `.*zt-logo.png`

}

package require http

set url http://max.demo.zentao.net/user-login-Lw==.html

set http [::http::geturl $url]
set html [::http::data $http]

regexp -- {<img src="(.*?)" .*>} $html match elem
puts "$elem"
