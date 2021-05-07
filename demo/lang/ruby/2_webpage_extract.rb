#!/usr/bin/env ruby

=begin

title=extract content from webpage
cid=0
pid=0

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `.*zt-logo.png`

=end

require "open-uri"

uri = 'http://max.demo.zentao.net/user-login-Lw==.html'
html = nil
open(uri) do |http|
  html = http.read
end

elem = html.match(/<img src="(.*?)" .*>/).captures
puts ""+elem[0]
