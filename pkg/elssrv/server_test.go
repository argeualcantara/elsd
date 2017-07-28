/**
 * (C) Copyright 2012-2016 HP Development Company, L.P.
 * Confidential computer software. Valid license from HP required for possession, use or copying.
 * Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
 * Computer Software Documentation, and Technical Data for Commercial Items are licensed
 * to the U.S. Government under vendor's standard commercial license.
 */
package elssrv

import (
	"flag"
	"github.com/hpcwp/elsd/pkg/api"
	"os"
	"testing"
)

var s ElsService
var service GRPCServer

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
	service = NewBasicService(RoutingKeyTableName, *dynamoDbAddr, *id, *secret, *token)
	s = NewService(RoutingKeyTableName, *dynamoDbAddr, *id, *secret, *token)

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
		request := api.AddRoutingKeyRequest{ServiceUri: "http://serviceA:80", RoutingKey: "555"}
		_, err := service.AddRoutingKey(nil, &request)
		if err != nil {
			t.Error("Failed adding routing key")
		}
	}

	{
		request := api.AddRoutingKeyRequest{ServiceUri: "http://serviceB:80", RoutingKey: "555"}
		_, err := service.AddRoutingKey(nil, &request)
		if err != nil {
			t.Error("Failed adding routing key")
		}
	}

	requestGet := api.RoutingKeyRequest{Id: "555"}
	response, err := service.ListServiceInstances(nil, &requestGet)

	if err != nil {
		t.Error("Failed listing routing key")
	}

	l := len(response.ServiceInstances)
	if l != 2 {
		t.Errorf("Expected 2 instances, got %d" , l)
	}

	for i := range response.ServiceInstances {
		t.Logf("Instance %v", response.ServiceInstances[i] )
	}


	{
		requestDel := api.DeleteRoutingKeyRequest{ServiceUri: "http://serviceA:80", RoutingKey: "555"}
		_, err3 := service.RemoveRoutingKey(nil, &requestDel)
		if err3 != nil {
			t.Error("Failed deleting routing key")
		}
	}


	{
		requestDel := api.DeleteRoutingKeyRequest{ServiceUri: "http://serviceB:80", RoutingKey: "555"}
		_, err3 := service.RemoveRoutingKey(nil, &requestDel)
		if err3 != nil {
			t.Error("Failed deleting routing key")
		}
	}



}

func TestGetKeys(t *testing.T) {
	request := api.AddRoutingKeyRequest{ServiceUri: "http://localhost:8072", RoutingKey: "666"}
	_, err := service.AddRoutingKey(nil, &request)
	if err != nil {
		t.Error("Failed adding routing key")
	}

	requestGet := api.RoutingKeyRequest{Id: "666"}
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
		request := api.AddRoutingKeyRequest{ServiceUri: "http://localhost:8072", RoutingKey: "100"}
		_, err := service.AddRoutingKey(nil, &request)
		if err != nil {
			t.Error("Failed adding routing key")
		}
	}

	{
		request := api.AddRoutingKeyRequest{ServiceUri: "http://localhost:8080", RoutingKey: "101"}
		_, err := service.AddRoutingKey(nil, &request)
		if err != nil {
			t.Error("Failed adding routing key")
		}
	}

	{
		requestGet := api.RoutingKeyRequest{Id: "100"}
		response, err2 := service.GetServiceInstanceByKey(nil, &requestGet)

		if err2 != nil {
			t.Error("Failed getting routing key")
		}

		if response.ServiceUri != "http://localhost:8072" {
			t.Error("Expected http://localhost:8072")
		}
	}

	{
		requestDel := api.DeleteRoutingKeyRequest{ServiceUri: "http://localhost:8072", RoutingKey: "100"}
		_, err3 := service.RemoveRoutingKey(nil, &requestDel)
		if err3 != nil {
			t.Error("Failed deleting routing key")
		}
	}

	{
		requestDel := api.DeleteRoutingKeyRequest{ServiceUri: "http://localhost:8080", RoutingKey: "101"}
		_, err3 := service.RemoveRoutingKey(nil, &requestDel)
		if err3 != nil {
			t.Error("Failed deleting routing key")
		}
	}

}

func TestGetKeysNegative(t *testing.T) {
	request := api.AddRoutingKeyRequest{ServiceUri: "http://localhost:8072", RoutingKey: "800"}
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
	request := api.AddRoutingKeyRequest{ServiceUri: "http://localhost:8072", RoutingKey: "900"}
	_, err := service.AddRoutingKey(nil, &request)
	if err != nil {
		t.Error("Failed adding routing key")
	}

	requestGet := api.RoutingKeyRequest{Id: "900"}
	_, err2 := service.GetServiceInstanceByKey(nil, &requestGet)

	if err2 != nil {
		t.Error("Failed getting routing key")
	}

	requestDel := api.DeleteRoutingKeyRequest{ServiceUri: "http://localhost:8072", RoutingKey: "900"}
	_, err3 := service.RemoveRoutingKey(nil, &requestDel)
	if err3 != nil {
		t.Error("Failed deleting routing key")
	}

	_, err4 := service.GetServiceInstanceByKey(nil, &requestGet)

	if err4 == nil {
		t.Error("Failed deleting routing key, the key is still there")
	}

}
