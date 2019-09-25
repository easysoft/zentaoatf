#!/usr/bin/env node
/**
[case]

title=check remote interface response
cid=0
pid=0

[group]
1. Send a request to interface http://xxx
2. Retrieve sessionID field from response json
3. Validate its format >> ^[a-z0-9]{26}

[esac]
*/

var http = require('http');

http.get('http://pms.zentao.net/?mode=getconfig', function(req) {
    let jsonStr = '';

    req.on('data', function(data) {
        jsonStr += data;
    });
    req.on('end', () => {
        if(req.statusCode === 200){
            try{
                var json = JSON.parse(jsonStr);
                console.log(">>" + json.sessionID)
            } catch(err){
                console.log('ERR: ' + err);
            }
        }
    });
});
