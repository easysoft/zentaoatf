#!/usr/bin/env node
/**
title=check remote interface response
cid=0
pid=0

1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Validate its format >> `^[a-z0-9]{26}`

*/

var http = require('http');

http.get('http://max.demo.zentao.net/pms/?mode=getconfig', function(req) {
    let jsonStr = '';

    req.on('data', function(data) {
        jsonStr += data;
    });
    req.on('end', () => {
        if(req.statusCode === 200){
            try{
                var json = JSON.parse(jsonStr);
                console.log(json.sessionID)
            } catch(err){
                console.log('ERR: ' + err);
            }
        }
    });
});
