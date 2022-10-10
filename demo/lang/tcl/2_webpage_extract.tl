#!/usr/bin/env tclsh
set case {

title=extract content from webpage
cid=0
pid=0

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `Bing images`

}

package require http
package require tls

http::register https 443 tls::socket

set url https://cn.bing.com/images/trending

set http [::http::geturl $url]
set html [::http::data $http]

regexp -- {<title>(.*?)</title>} $html match title
puts "$title"
