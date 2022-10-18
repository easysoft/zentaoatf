#!/usr/bin/env ruby
=begin

title=check remote interface response
cid=0
pid=0

1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Check its format >> `^[0-9]{8}`

=end

require 'uri'
require 'net/http'
require 'json'

uri = URI('https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1')
res = Net::HTTP.get_response(uri)

json = JSON.parse(res.body)

puts json['images'][0]['startdate']

