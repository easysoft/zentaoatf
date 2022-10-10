#!/usr/bin/env node
/**
title=check remote interface response
cid=0
pid=0

1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Validate its format >> `^[a-z0-9]{8}`

*/

const https = require('https');

https.get('https://cn.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1', function(req) {
    let jsonStr = '';

    req.on('data', function(data) {
        jsonStr += data;
    });
    req.on('end', () => {
        if(req.statusCode === 200){
            try{
                const json = JSON.parse(jsonStr);
                console.log(json.images[0].startdate)
            } catch(err){
                console.log('ERR: ' + err);
            }
        }
    });
});
