package elscli

import (
	"context"
	"github.com/hpcwp/elsd/pkg/api"
	"github.com/prometheus/common/log"
)

func GetServiceInstanceByKey(client api.ElsClient, routingKey string) (*api.ServiceInstanceReponse, error) {
	req := &api.RoutingKeyRequest{routingKey}
	resp, err := client.GetServiceInstanceByKey(context.Background(), req)
	if err != nil {
		log.Fatalf("Error gettting routing key", err)
		return nil, err
	}

	log.Info("Rotung key %s and tags %s", resp.GetServiceUri(), resp.GetTags())
	return resp, nil
}

func AddServiceInstance(client api.ElsClient, routingKey string, uri string, tags []string) (*api.ServiceInstanceReponse, error) {
	req := &api.AddRoutingKeyRequest{ uri, tags[0], routingKey}
	resp, err := client.AddRoutingKey(context.Background(),req)
	if err !=nil {
		log.Fatalf("Error adding service instanve", err)
		return nil, err
	}

	log.Info("ServiceInstnace added", resp.GetServiceUri(), resp.GetTags())
	return resp, nil

}
