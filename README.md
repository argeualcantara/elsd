# ELSd

Entity Locator Service

## Building

```
docker-compose build
```

### Modifying the code

Generating gRPC client and server interfaces.

```shell
protoc -I pkg/api/ pkg/api/els.proto --go_out=plugins=grpc:pkg/api
```

Updating dependencies.

```shell
dep ensure -update
```

To update a dependency to a new version, you might run

```shell
dep ensure github.com/pkg/errors@^0.8.0
```

## Running

```shell
docker-compose up
```

## Usage

You can build the client, add some routing keys and get them

```shell
go install  github.com/hpcwp/elsd/cmd/elscli
```

Test elsd server is working

```shell
elscli

```

```shell
elscli -method Add   123 http://localhost:8072 rw
elscli -method Add   123 http://localhost:8080 r
elscli -method Add   124 http://localhost:8072 rw
elscli -method Add   125 http://localhost:8080 rw
```

```shell
elscli -method Get   125
```

```shell
elscli -method Remove  125 http://localhost:8080
```

Remove is an idempotent operation, so you can call it multiple times with the same result,
also if a key does not exist it will not result in any error

## Testing

To run integration testing, start the dynamodb container and run the tests

```shell
docker-compose up dynamodb

```

```shell
go test $(go list ./... | grep -v /vendor/)
```

## Prometheus Metrics

```http
http://localhost:8080/metrics
```

## Using AWS CLI

```shell
export AWS_ACCESS_KEY_ID=123
export AWS_SECRET_ACCESS_KEY=123
aws dynamodb list-tables --endpoint-url http://localhost:8000
 ```

## Configuring DynamoDB on AWS

To use the commands bellow your aws cli version must be >= `1.14.69`

### us-west-2

```bash

aws dynamodb create-table \
    --table-name RoutingKeys.v1 \
    --attribute-definitions AttributeName=Id,AttributeType=S AttributeName=Uri,AttributeType=S \
    --key-schema AttributeName=Id,KeyType=HASH AttributeName=Uri,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --stream-specification StreamEnabled=true,StreamViewType=NEW_AND_OLD_IMAGES \
    --region=us-west-2

```
### eu-west-1
````bash

aws dynamodb create-table \
    --table-name RoutingKeys.v1 \
    --attribute-definitions AttributeName=Id,AttributeType=S AttributeName=Uri,AttributeType=S \
    --key-schema AttributeName=Id,KeyType=HASH AttributeName=Uri,KeyType=RANGE \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
    --stream-specification StreamEnabled=true,StreamViewType=NEW_AND_OLD_IMAGES \
    --region=eu-west-1

````
### us-west-2 (GLOBAL TABLE)
```bash

aws dynamodb create-global-table \
--global-table-name RoutingKeys.v1 \
--replication-group RegionName=us-west-2 RegionName=eu-west-1 \
--region us-west-2

```