package elssrv

// This file contains the Service definition, and a basic service
// implementation. It also includes service middlewares.

import (
	"errors"
	"github.com/galo/els-go/pkg/api"
	"golang.org/x/net/context"
	"github.com/hpcwp/els-go/dynamodb/routingkeys"
)

// Service describes a service that adds things together.
type ElsService interface {
	GetServiceInstanceByKey(ctx context.Context, routingKey *api.RoutingKey) (*api.ServiceInstance, error)
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
func (bs basicElsService) GetServiceInstanceByKey(ctx context.Context, routingKey *api.RoutingKey) (*api.ServiceInstance, error) {

	if routingKey.Id == "" {
		return &api.ServiceInstance{}, ErrInvalid
	}

	serviceInstance := bs.rksrv.Get(routingKey.Id)


	if serviceInstance == nil {
		return nil, ErrNotFound
	}
	if len(serviceInstance.Stacks) == 0 {
		return nil, ErrNotFound
	}

	// We just return teh first service url
	serviceUrl := serviceInstance.Stacks[0].Name
	if serviceUrl == nil {
		return nil, ErrNotFound
	}

	srvInstance := api.ServiceInstance{*serviceUrl, "rw"}
	return &srvInstance, nil
}

const routingKeyTableName  = "routingKeys"

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() ElsService {
	rk := routingkeys.New(routingKeyTableName)

	return basicElsService{rk}
}

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(ElsService) ElsService

// ErrEmpty is returned when input is invalid
var ErrInvalid = errors.New("Invalid routing key")
