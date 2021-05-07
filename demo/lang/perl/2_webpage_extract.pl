#!/usr/bin/env perl
=pod
title=extract content from webpage
cid=0
pid=0

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `.*zt-logo.png`

=cut

use LWP::Simple; # need LWP::Simple module
$html = get('http://max.demo.zentao.net/user-login-Lw==.html');

if ( $html =~ /<img src="(.*?)" .*>/ ) {
  print "$1\n";
}
