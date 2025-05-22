# ok200

## create certificate for https
```bash
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=*" -addext "subjectAltName = DNS:*,IP:0.0.0.0"
```

## run server
```bash
go run main.go -p 8888 -c utf8 -https -cert cert.pem -key key.pem
```



