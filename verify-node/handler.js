"use strict"
const http = require('http');

module.exports = (context, callback) => {
    var data = "";

    http.get("http://gateway:8080/certificates/callback", (res) => {
        res.setEncoding('utf8');
        res.on('data', function(chunk) {
            data += chunk;
        });

        res.on('end', () => {
            let json = JSON.parse(data)
            verify(json.pem)
        });
    }).on('error', (e) => {
        console.error(`Got error: ${e.message}`);
        callback(undefined, {status: "error", error: e.message});
    });

    function verify(publicKey) {
        // if request is not signed treat it as untrusted
        let auth = process.env.Http_Authorization;
        if (auth === undefined) {
            callback(undefined, {status: 403, pem: publicKey});
            return
        }

        console.log('authorization:' + auth);
        console.log("public key:\n" + publicKey);
        console.log("(request-target): %s %s\n", strings.ToLower(r.Method), r.URL.Path)
        console.log("date: %s\n", headerOrDefault(r, "Date", ""))
        console.log("content-type: %s\n", headerOrDefault(r, "Content-Type", ""))
        console.log("digest: %s\n", headerOrDefault(r, "Digest", ""))
        console.log("body:\n%s\n", string(body))
        console.log("authorization: %s\n", headerOrDefault(r , "Authorization", ""))


        callback(undefined, {status: status, pem: publicKey});
    }
}
