#!/usr/bin/env ruby
=begin
[case]

title=check remote interface response
cid=0
pid=0

[group]
1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Check its format >> ^[a-z0-9]{26}

[esac]
=end

require "open-uri"
require "json"

uri = 'http://pms.zentao.net/?mode=getconfig'
html = nil
open(uri) do |http|
  html = http.read
end

json = JSON.parse(html)   # need to install json library,  gem install json
puts '>> ' + json['sessionID']
