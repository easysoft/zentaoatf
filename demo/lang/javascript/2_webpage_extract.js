#!/usr/bin/env node
/**
[case]
title=extract content from webpage
cid=0
pid=0

[group]
  1. Load web page from url http://xxx 
  2. Retrieve img element zt-logo.png in html 
  3. Check img exist >> .*zt-logo.png

[esac]
*/

var http = require('http');

http.get('http://pms.zentao.net/user-login.html', function(req) {
    let html = '', image = '';

    req.on('data', function(data) {
        html += data;
    });
    req.on('end', () => {
        var res = html.match(/<img\ssrc='(.*?)'/);
        if (res.length > 1) {
            image = res[1]
        }

        console.log(">>" + res[1])
    });
});