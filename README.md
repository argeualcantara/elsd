# ELSd

Entity Locator Service

## Building

```
$ docker-compose build
```

## Running

```
$ docker-compose up
```

## Testing

```
go run cmd/elscli/main.go  -grpc.addr localhost:8082 -method GetServiceInstanceByKey mykey
```


```
$ elscli -grpc.addr localhost:8082 -method GetServiceInstanceByKey mykey
```

