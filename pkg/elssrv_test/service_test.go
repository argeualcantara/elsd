package elssrv_test

import (
	"testing"
	"os"
	"flag"
	"github.com/hpcwp/elsd/pkg/elssrv"
	"github.com/hpcwp/elsd/pkg/api"
)

var service elssrv.ElsService

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags

	var (
		dynamoDbAddr = flag.String("dynamodb.addr", "http://localhost:8000", "DynamoDb (HTTP) address")
		id           = flag.String("aws.id", "123", "AWS id")
		secret       = flag.String("aws.secret", "123", "AWS secret")
		token        = flag.String("aws.token", "", "AWS token")
	)
	flag.Parse()

	// Business domain.
	service = elssrv.NewBasicService(elssrv.RoutingKeyTableName, *dynamoDbAddr, *id, *secret, *token)


	os.Exit(m.Run())
}

func TestAddKeys(t *testing.T) {
	request := api.AddRoutingKeyRequest{ServiceUri: "http://localhost:8072", RoutingKey: "123"}
	_, err := service.AddRoutingKey(nil, &request)
	if err != nil {
		t.Error("Failed adding routing key")
	}

}

func TestGetKeys(t *testing.T) {
	request := api.AddRoutingKeyRequest{ServiceUri: "http://localhost:8072", RoutingKey: "123"}
	_, err := service.AddRoutingKey(nil, &request)
	if err != nil {
		t.Error("Failed adding routing key")
	}


	requestGet := api.RoutingKeyRequest{Id: "123"}
	response, err2 := service.GetServiceInstanceByKey(nil, &requestGet)

	if err2 != nil {
		t.Error("Failed getting routing key")
	}

	if response.ServiceUri != "http://localhost:8072" {
		t.Error("Expected http://localhost:8072")
	}

}

func TestGetKeysNegative(t *testing.T) {
	request := api.AddRoutingKeyRequest{ServiceUri: "http://localhost:8072", RoutingKey: "123"}
	_, err := service.AddRoutingKey(nil, &request)
	if err != nil {
		t.Error("Failed adding routing key")
	}


	requestGet := api.RoutingKeyRequest{Id: "124"}
	_, err2 := service.GetServiceInstanceByKey(nil, &requestGet)

	if err2 == nil {
		t.Error("Failed getting routing key")
	}


}

