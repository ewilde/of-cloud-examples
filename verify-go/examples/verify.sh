#!/usr/bin/env bash
# verify a signature
# ./verify.sh ./sigining_string.txt  ./signature_base64 ./public_key.pem

signing_string=$1
signature=$2
publickey=$3

cat $signature | base64 -D > /tmp/signature.sha256

openssl dgst -sha256 -verify $publickey -signature /tmp/signature.sha256 $signing_string
