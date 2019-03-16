'use strict'
const http = require('http');
const crypto = require('crypto');
const hash = crypto.createHash('sha256');

module.exports = (body, callback) => {
    var data = '';

    http.get('http://gateway:8080/certificates/callback', (res) => {
        res.setEncoding('utf8');
        res.on('data', function(chunk) {
            data += chunk;
        });

        res.on('end', () => {
            let json = JSON.parse(data);
            verify(json.pem);
        });
    }).on('error', (e) => {
        console.warn(`Got error: ${e.message}`);
        callback(undefined, {status: 'error', error: e.message});
    });

    function verify(publicKey) {
        let auth = process.env.Http_Authorization;
        if (auth === undefined) {
            callback(undefined, {status: 403, message: 'missing authorization header'});
            return
        }

        let signingString = getSigningString(body);
        let signature = getSignature(auth);
        console.warn(sprintf('signing string:\n%s', signingString));
        console.warn('public key:\n' + publicKey);
        console.warn('authorization: %s\n', auth);
        console.warn('signature: %s\n', signature);

        var verify = crypto.createVerify('RSA-SHA256');
        verify.update(signingString);

        var verification = verify.verify(publicKey, signature, 'base64');

        if (verification) {
            console.warn("Verification OK");
            callback(undefined, {status: '200', verification: verification});

        } else {
            console.warn("Verification failed");
            callback(undefined, {status: '401', verification: verification});
        }
    }

    function getSigningString(body) {

        hash.update(body);
        let digest = hash.digest('base64');

        let res = '';
        res += sprintf('%s\n', getRequestTarget());
        res += sprintf('host: %s\n', process.env.Http_X_Forwarded_Host);
        res += sprintf('date: %s\n', process.env.Http_Date);
        res += sprintf('content-type: %s\n', process.env.Http_Content_Type);
        res += sprintf('digest: SHA-256=%s\n', digest);
        res += sprintf('content-length: %s', body.length);

        return res
    }

    function getSignature(authHeader) {
        return authHeader.split(',')[3].split('="')[1] + '='
    }

    function getRequestTarget() {
        let path = process.env.Http_Path;
        if (path === '/') {
            path = '';
        }

        let name = process.env.Http_Host.split(':')[0];

        let requestTarget = '(request-target): ' + process.env.Http_Method.toLowerCase() + ' /function/' + name + path;
        let query = process.env.Http_Query;
        if (query !== undefined) {
            requestTarget = requestTarget + '?' + query
        }

        return requestTarget
    }

    function sprintf(str) {
        var args = [].slice.call(arguments, 1),
            i = 0;

        return str.replace(/%s/g, () => args[i++]);
    }
};
