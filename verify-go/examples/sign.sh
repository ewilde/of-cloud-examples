#!/bin/bash
# Create an authorization headers based on the [http signatures draft](https://tools.ietf.org/html/draft-cavage-http-signatures-10):
#
#
# Usage: sign url_path host date content_type body private_key
#
# example: sign.sh 'POST /foo?param=value&pet=dog' 'example.com' 'Sun, 05 Jan 2014 21:31:40 GMT' 'application/json' '{"hello": "world"}' private_key.pem
#

path=$(echo $1 |  tr '[:upper:]' '[:lower:]')
host=$2
date=$3
content_type=$4
body=$5
length=${#body}
private_key=$6

digest=$(echo -n ${body}| openssl dgst -sha256 -binary | base64)
headers="(request-target): ${path}\nhost: ${host}\ndate: ${date}\ncontent-type: ${content_type}\ndigest: SHA-256=${digest}\ncontent-length: ${length}"

signature=$(echo -en ${headers} | openssl dgst -sha256 -sign ${private_key} | base64)
echo digest: ${digest}
echo signature: ${signature}
echo Authorization: Signature keyId=\"Test\",algorithm=\"rsa-sha256\",headers=\"\(request-target\) host date content-type digest content-length\",signature=\"${signature}\"
