## Description
This example function takes a CSR and creates a self-signed certificate


### Example

```bash
curl http://192.168.99.100:31112/function/cert-sign-go \
 -d '{
        "Host": "example.com",
        "RSAKeySize": 2048,
        "ValidFor": 63072000000000000
     }'
```