"use strict"
const http = require('http');


module.exports = (context, callback) => {
    var data = '';

    http.get('http://ip.jsontest.com/', (res) => {
        res.setEncoding('utf8');
        res.on('data', function(chunk) {
            data += chunk;
        });

        res.on('end', () => {
            callback(undefined, {status: 'success', ip: data});
        });
    }).on('error', (e) => {
        console.warn(`Got error: ${e.message}`);
        callback(undefined, {status: 'error', error: e.message});
    });
}
