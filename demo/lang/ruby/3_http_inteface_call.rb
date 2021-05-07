#!/usr/bin/env ruby
=begin

title=check remote interface response
cid=0
pid=0

1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Check its format >> `^[a-z0-9]{26}`

=end

require "open-uri"
require "json"

uri = 'http://max.demo.zentao.net/pms/?mode=getconfig'
html = nil
open(uri) do |http|
  html = http.read
end

json = JSON.parse(html)   # need json library (gem install json)
puts json['sessionID']
