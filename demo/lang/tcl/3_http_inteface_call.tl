#!/usr/bin/env tclsh
set case {
[case]

title=check remote interface response
cid=0
pid=0

[group]
1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Check its format >> ^[a-z0-9]{26}

[esac]
}

package require http
package require json

set url http://pms.zentao.net/?mode=getconfig

set http [::http::geturl $url]
set jsonStr [::http::data $http]

# need to install json library, you may use ActiveTcl
set jsonObj [json::json2dict $jsonStr]
puts ">> [dict get $jsonObj sessionID]"
