# LitePub Relay

## Setup
### Generate RSA Key
```
openssl genrsa -out privatekey.pem 4096
openssl rsa -in privatekey.pem -pubout -out publickey.pem
```

