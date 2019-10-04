#!/usr/bin/env perl
=pod
[case]
title=extract content from webpage
cid=0
pid=0

[group]
  1. Load web page from url http://xxx 
  2. Retrieve img element zt-logo.png in html 
  3. Check img exist >> .*zt-logo.png

[esac]
=cut

use LWP::Simple; # need LWP::Simple module
$html = get('http://pms.zentao.net/user-login.html');

if ( $html =~ /<img src='(.*?)' .*>/ ) {
  print ">> $1\n";
}