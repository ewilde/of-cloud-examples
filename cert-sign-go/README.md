## Description
This example function takes a CSR and creates a self-signed certificate


### Example

```bash
curl https://ewilde.o6s.io/cert-sign-go/function/cert-sign-go \
 -d '{
        "Host": "example.com",
        "RSAKeySize": 2048,
        "ValidFor": 63072000000000000
     }'
```