#!/usr/bin/env tclsh
set case {

title=check remote interface response
cid=0
pid=0

1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Check its format >> `^[a-z0-9]{8}`

}

package require http
package require tls
package require json

http::register https 443 tls::socket

set url https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1

set http [::http::geturl $url]
set jsonStr [::http::data $http]

regexp -- {"startdate":"(.*?)"} $jsonStr match date
puts "$date"

# set jsonObj [json::json2dict $jsonStr]
# set images [dict get $jsonObj images]
# puts [dict get $images startdate]
