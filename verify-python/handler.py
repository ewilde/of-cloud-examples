import base64
import os
import hashlib
import logging
import requests

def handle(req):

    """handle a request to the function
    Args:
        req (str): request body
    """

    body_digest = get_digest(req)
    signing_string = get_signing_string(req, body_digest)
    signature = get_signature()
    public_key = get_public_key()

    logging.warn('\n' + signing_string)
    logging.warn('\n' + public_key)

    if verify(signing_string, public_key, signature):
        logging.warn('{status: "200", verification: true}')
        return '{status: "200", verification: true}'
    else:
        logging.warn('{status: "200", verification: false}')
        return '{status: "200", verification: false}'


def verify(signing_string, public_key_pem, signature):
    from Crypto.PublicKey import RSA
    from Crypto.Signature import PKCS1_v1_5
    from Crypto.Hash import SHA256
    from base64 import b64decode

    public_key = RSA.importKey(public_key_pem)
    signer = PKCS1_v1_5.new(public_key)
    digest = SHA256.new()
    digest.update(signing_string)

    if signer.verify(digest, b64decode(signature)):
        return True
    return False


def get_signature():
    auth_header = os.environ['Http_Authorization']
    signature = auth_header.split(',')[3].split('="')[1] + '='

    return signature


def get_digest(body):
    sha_hash = hashlib.sha256()
    sha_hash.update(body)
    digest = base64.b64encode(sha_hash.digest()).decode()

    return digest


def get_signing_string(body, digest):

    res = ''
    res += '{}\n'.format(get_request_target())
    res += 'host: {}\n'.format(os.environ['Http_X_Forwarded_Host'])
    res += 'date: {}\n'.format(os.environ['Http_Date'])
    res += 'content-type: {}\n'.format(os.environ['Http_Content_Type'])
    res += 'digest: SHA-256={}\n'.format(digest)
    res += 'content-length: {}'.format(len(body))
    
    return res


def get_request_target():

    path = os.environ['Http_Path']
    if path == '/':
        path = ''

    name = os.environ['Http_Host'].split(':')[0]

    request_target = '(request-target): ' + os.environ['Http_Method'].lower() + ' /function/' + name + path

    if 'Http_Query' in os.environ:
        query = os.environ['Http_Query']
        request_target = request_target + '?' + query

    return request_target


def get_public_key():

    req = requests.get('http://gateway:8080/certificates/callback')

    return req.json()['pem']
