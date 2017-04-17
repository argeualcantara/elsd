# ELSd

Entity Locator Service

## Building

```
$ docker-compose build
```

### Modifying the code 

Generating gRPC client and server interfaces.  

```
$ protoc -I pkg/api/ pkg/api/els.proto --go_out=plugins=grpc:pkg/api
```

Updating dependencies.  
```
$ dep ensure -update
```

To update a dependency to a new version, you might run

```
$ dep ensure github.com/pkg/errors@^0.8.0

```

## Running

```
$ docker-compose up
```

## Testing

```
$ go run cmd/elscli/main.go  -grpc.addr localhost:8082 -method GetServiceInstanceByKey mykey
```


```
$ elscli -grpc.addr localhost:8082 -method GetServiceInstanceByKey mykey
```

