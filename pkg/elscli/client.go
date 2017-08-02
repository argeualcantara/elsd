/**
 * (C) Copyright 2012-2016 HP Development Company, L.P.
 * Confidential computer software. Valid license from HP required for possession, use or copying.
 * Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
 * Computer Software Documentation, and Technical Data for Commercial Items are licensed
 * to the U.S. Government under vendor's standard commercial license.
 */
package elscli

import (
	"context"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/hpcwp/elsd/pkg/api"
	"github.com/prometheus/common/log"
)

func GetServiceInstanceByKey(client api.ElsClient, routingKey string) (*api.ServiceInstanceResponse, error) {
	req := &api.RoutingKeyRequest{routingKey}
	resp, err := client.GetServiceInstanceByKey(context.Background(), req)
	if err != nil {
		log.Fatalf("Error getting routing key", err)
		return nil, err
	}

	log.Info("Roting key %s and tags %s", resp.GetServiceUri(), resp.GetTags())
	return resp, nil
}

func ListServiceInstances(client api.ElsClient, routingKey string) (*api.ServiceInstanceListResponse, error) {
	req := &api.RoutingKeyRequest{routingKey}
	resp, err := client.ListServiceInstances(context.Background(), req)
	if err != nil {
		log.Fatalf("Error listing routing key", err)
		return nil, err
	}

	for i := range resp.ServiceInstances {
		log.Info("Roting key %s and tags %s", resp.ServiceInstances[i].GetServiceUri(), resp.ServiceInstances[i].GetTags())
	}
	return resp, nil
}

func AddServiceInstance(client api.ElsClient, routingKey string, uri string, tags []string) (*api.ServiceInstanceResponse, error) {
	req := &api.AddRoutingKeyRequest{uri, tags[0], routingKey}
	resp, err := client.AddRoutingKey(context.Background(), req)
	if err != nil {
		log.Fatalf("Error adding service instance", err)
		return nil, err
	}

	log.Info("ServiceInstnace added", resp.GetServiceUri(), resp.GetTags())
	return resp, nil

}

func RemoveServiceInstance(client api.ElsClient, routingKey string, uri string) (*google_protobuf.Empty, error) {
	req := &api.DeleteRoutingKeyRequest{uri, routingKey}
	resp, err := client.RemoveRoutingKey(context.Background(), req)
	if err != nil {
		log.Fatalf("Error deleting routing key ", err)
		return nil, err
	}
	log.Info("RoutingKey deleted")
	return resp, nil
}
