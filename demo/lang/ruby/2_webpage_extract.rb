#!/usr/bin/env ruby
=begin
[case]
title=extract content from webpage
cid=0
pid=0

[group]
  1. Load web page from url http://xxx 
  2. Retrieve img element zt-logo.png in html 
  3. Check img exist >> `.*zt-logo.png`

[esac]
=end

require "open-uri"

uri = 'http://pms.zentao.net/user-login.html'
html = nil
open(uri) do |http|
  html = http.read
end

elem = html.match(/<img src='(.*?)' .*>/).captures
puts '>> ' + elem[0]