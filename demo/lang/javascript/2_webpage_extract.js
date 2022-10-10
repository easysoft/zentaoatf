#!/usr/bin/env node
/**
title=extract content from webpage
cid=0
pid=0

1. Load web page from url http://xxx
2. Retrieve img element zt-logo.png in html
3. Check img exist >> `必应`

*/

const https = require('https');

https.get('https://cn.bing.com', function(req) {
    let html = ''

    req.on('data', function(data) {
        html += data;
    });
    req.on('end', () => {
        const res = html.match(/<title>(.*?)</);
        if (res.length > 1) {
            console.log(res[1])
        }
    });
});
