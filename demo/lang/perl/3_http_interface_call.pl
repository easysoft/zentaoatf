#!/usr/bin/env perl
=pod

title=check remote interface response
cid=0
pid=0

1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Check its format >> `^[a-z0-9]{26}`

=cut

use LWP::Simple; # need LWP::Simple module
$json = get('http://max.demo.zentao.net/pms/?mode=getconfig');

if ( $json =~ /"sessionID":"([^"]*)"/ ) {
  print "$1\n";
}
