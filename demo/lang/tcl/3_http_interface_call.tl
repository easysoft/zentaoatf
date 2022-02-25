#!/usr/bin/env tclsh
set case {

title=check remote interface response
cid=0
pid=0

  1. Send a request to interface http://xxx
  2. Retrieve sessionID field from response json
  3. Check its format >> `^[a-z0-9]{26}`

}

package require http
package require json

set url http://max.demo.zentao.net/pms/?mode=getconfig

set http [::http::geturl $url]
set jsonStr [::http::data $http]

# need json library, you may use ActiveTcl
set jsonObj [json::json2dict $jsonStr]
puts "[dict get $jsonObj sessionID]"
