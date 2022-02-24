#!/usr/bin/env node
/**
title=extract content from webpage
cid=0
pid=0

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `.*.png`

*/

var http = require('http');

http.get('http://max.demo.zentao.net/user-login-Lw==.html', function(req) {
    let html = '', image = '';

    req.on('data', function(data) {
        html += data;
    });
    req.on('end', () => {
        var res = html.match(/<img\ssrc="(.*?)"/);
        if (res.length > 1) {
            image = res[1]
            console.log(image)
        }
    });
});
