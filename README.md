# Nats Proxy

Simple Golang HTTP application to proxy rpc requests to nats server

## 1. Run with Docker

1. **Build**

```shell script
make build
docker build . -t api-rest
```

2. **Run**

```shell script
docker run -p 3000:3000 api-rest
```
