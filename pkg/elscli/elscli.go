package elscli

import (
	"github.com/galo/els-go/pkg/api"
	"context"
	"github.com/prometheus/common/log"
)




func GetServiceInstanceByKey(client api.ElsClient, routingKey string)  (*api.ServiceInstance, error){
	req := &api.RoutingKey{routingKey}
	resp,err :=client.GetServiceInstanceByKey(context.Background(), req)
	if err !=nil {
		log.Fatalf("Error gettting routing jey: %v", err)
		return nil,err
	}

	log.Info("Rotung key %s and tags %s", resp.GetServiceUri(), resp.GetTags())
	return resp,nil
}
