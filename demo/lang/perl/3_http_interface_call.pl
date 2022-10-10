#!/usr/bin/env perl
=pod

title=check remote interface response
cid=0
pid=0

1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Check its format >> `^[a-z0-9]{8}`

=cut

use LWP::Simple; # need LWP::Simple module
$json = get('https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1');

if ( $json =~ /"startdate":"([^"]*)"/ ) {
  print "$1\n";
}
