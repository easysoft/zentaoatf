#!/usr/bin/env ruby

=begin

title=extract content from webpage
cid=0
pid=0

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `必应`

=end

require 'uri'
require 'net/http'

uri = URI('https://cn.bing.com')
res = Net::HTTP.get_response(uri)
html = res.body

elem = html.match(/<title>(.*?)</).captures
puts ""+elem[0]