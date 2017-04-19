package elssrv

// This file contains the Service definition, and a basic service
// implementation. It also includes service middlewares.

import (
	"errors"
	"github.com/hpcwp/elsd/pkg/api"
	"github.com/hpcwp/elsd/pkg/dynamodb/routingkeys"
	"golang.org/x/net/context"
)

// Service describes a service that adds things together.
type ElsService interface {
	GetServiceInstanceByKey(ctx context.Context, request *api.RoutingKeyRequest) (*api.ServiceInstanceReponse, error)

	// Add a routingKey to a service
	AddRoutingKey(context.Context, *api.AddRoutingKeyRequest) (*api.ServiceInstanceReponse, error)
}

type ServiceInstance struct {
	Url      string `json:"url"`
	Metadata string `json:"metadata"`
}

type basicElsService struct{
	rksrv *routingkeys.Service
}

// Errors
var (
	ErrNotFound = errors.New("ServiceInstance not found ")
)

// The implementation of the service
func (bs basicElsService) GetServiceInstanceByKey(ctx context.Context, routingKey *api.RoutingKeyRequest) (*api.ServiceInstanceReponse, error) {

	if routingKey.Id == "" {
		return &api.ServiceInstanceReponse{}, ErrInvalid
	}

	serviceInstance := bs.rksrv.Get(routingKey.Id)


	if serviceInstance == nil {
		return nil, ErrNotFound
	}
	if len(serviceInstance.ServiceInstances) == 0 {
		return nil, ErrNotFound
	}

	// We just return the first service url
	serviceUrl := serviceInstance.ServiceInstances[0].Uri
	if serviceUrl == "" {
		return nil, ErrNotFound
	}

	srvInstance := api.ServiceInstanceReponse{serviceUrl, "rw"}
	return &srvInstance, nil
}

// The implementation of teh service
func (bs basicElsService) AddRoutingKey(ctx context.Context, addRoutingKeyRequest *api.AddRoutingKeyRequest) (*api.ServiceInstanceReponse, error) {
	if addRoutingKeyRequest.ServiceUri== "" {
		return &api.ServiceInstanceReponse{}, ErrInvalid
	}
	if addRoutingKeyRequest.RoutingKey== "" {
		return &api.ServiceInstanceReponse{}, ErrNotFound
	}

	instance := &routingkeys.ServiceInstance{addRoutingKeyRequest.RoutingKey,
		addRoutingKeyRequest.ServiceUri,
		[]string{addRoutingKeyRequest.Tags}}

	bs.rksrv.Add(instance)

	return &api.ServiceInstanceReponse{instance.Uri,
		instance.Tags[0]}, nil

}


const RoutingKeyTableName  = "routingKeys"


// NewBasicService returns a naïve dynamoDb implementation of Service.
func NewBasicService(tableName string, dynamoAddr string, id string , secret string , token string) ElsService {
	rk := routingkeys.New(tableName, dynamoAddr, id, secret, token)

	return basicElsService{rk}
}

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(ElsService) ElsService

// ErrEmpty is returned when input is invalid
var ErrInvalid = errors.New("Invalid routing key")
