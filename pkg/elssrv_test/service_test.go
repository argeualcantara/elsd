package elssrv_test

import (
	"flag"
	"github.com/hpcwp/elsd/pkg/api"
	"github.com/hpcwp/elsd/pkg/elssrv"
	"os"
	"testing"
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

func TestListKeys(t *testing.T) {
	{
		request := api.AddRoutingKeyRequest{ServiceUri: "http://serviceA:80", RoutingKey: "123"}
		_, err := service.AddRoutingKey(nil, &request)
		if err != nil {
			t.Error("Failed adding routing key")
		}
	}

	{
		request := api.AddRoutingKeyRequest{ServiceUri: "http://serviceB:80", RoutingKey: "123"}
		_, err := service.AddRoutingKey(nil, &request)
		if err != nil {
			t.Error("Failed adding routing key")
		}
	}

	requestGet := api.RoutingKeyRequest{Id: "123"}
	response, err := service.ListServiceInstances(nil, &requestGet)

	if err != nil {
		t.Error("Failed listing routing key")
	}

	if len(response.ServiceInstances) != 2 {
		t.Error("Expected two instances")

	}

	if response.ServiceInstances[0].ServiceUri != "http://serviceA:80" {
		t.Error("Expected http://localhost:8072")
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

func TestGetManyKeys(t *testing.T) {
	{
		request := api.AddRoutingKeyRequest{ServiceUri: "http://localhost:8072", RoutingKey: "123"}
		_, err := service.AddRoutingKey(nil, &request)
		if err != nil {
			t.Error("Failed adding routing key")
		}
	}

	{
		request := api.AddRoutingKeyRequest{ServiceUri: "http://localhost:8080", RoutingKey: "125"}
		_, err := service.AddRoutingKey(nil, &request)
		if err != nil {
			t.Error("Failed adding routing key")
		}
	}

	{
		requestGet := api.RoutingKeyRequest{Id: "123"}
		response, err2 := service.GetServiceInstanceByKey(nil, &requestGet)

		if err2 != nil {
			t.Error("Failed getting routing key")
		}

		if response.ServiceUri != "http://localhost:8072" {
			t.Error("Expected http://localhost:8072")
		}
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

func TestDeleteKeys(t *testing.T) {
	request := api.AddRoutingKeyRequest{ServiceUri: "http://localhost:8072", RoutingKey: "123"}
	_, err := service.AddRoutingKey(nil, &request)
	if err != nil {
		t.Error("Failed adding routing key")
	}

	requestGet := api.RoutingKeyRequest{Id: "123"}
	_, err2 := service.GetServiceInstanceByKey(nil, &requestGet)

	if err2 != nil {
		t.Error("Failed getting routing key")
	}

	requestDel := api.DeleteRoutingKeyRequest{ServiceUri: "http://localhost:8072", RoutingKey: "123"}
	_, err3 := service.RemoveRoutingKey(nil, &requestDel)
	if err3 != nil {
		t.Error("Failed deleting routing key")
	}

	_, err4 := service.GetServiceInstanceByKey(nil, &requestGet)

	if err4 == nil {
		t.Error("Failed deleting routing key, the key is still there")
	}

}
