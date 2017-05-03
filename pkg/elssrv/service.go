/**
 * (C) Copyright 2012-2016 HP Development Company, L.P.
 * Confidential computer software. Valid license from HP required for possession, use or copying.
 * Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
 * Computer Software Documentation, and Technical Data for Commercial Items are licensed
 * to the U.S. Government under vendor's standard commercial license.
 */
package elssrv

// This file contains the Service definition, and a basic service
// implementation. It also includes service middlewares.

import (
	"errors"
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/hpcwp/elsd/pkg/api"
	"github.com/hpcwp/elsd/pkg/dynamodb/routingkeys"
	"golang.org/x/net/context"
)

// Service describes a service that adds things together.
type ElsService interface {
	GetServiceInstanceByKey(ctx context.Context, request *api.RoutingKeyRequest) (*api.ServiceInstanceResponse, error)

	// Add a routingKey to a service
	AddRoutingKey(context.Context, *api.AddRoutingKeyRequest) (*api.ServiceInstanceResponse, error)

	// Delete a routingKey to a service
	RemoveRoutingKey(context.Context, *api.DeleteRoutingKeyRequest) (*google_protobuf.Empty, error)
}

type ServiceInstance struct {
	Url      string `json:"url"`
	Metadata string `json:"metadata"`
}

type basicElsService struct {
	rksrv *routingkeys.Service
}

// Errors
var (
	// ErrEmpty is returned when input is invalid
	ErrInvalid  = errors.New("invalid routing key")
	ErrNotFound = errors.New("service instance not found ")
)

// The implementation of the service
func (bs basicElsService) GetServiceInstanceByKey(ctx context.Context, routingKey *api.RoutingKeyRequest) (*api.ServiceInstanceResponse, error) {

	if routingKey.Id == "" {
		return &api.ServiceInstanceResponse{}, ErrInvalid
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

	srvInstance := api.ServiceInstanceResponse{serviceUrl, "rw"}
	return &srvInstance, nil
}

// The implementation of the service
func (bs basicElsService) AddRoutingKey(ctx context.Context, addRoutingKeyRequest *api.AddRoutingKeyRequest) (*api.ServiceInstanceResponse, error) {
	if addRoutingKeyRequest.ServiceUri == "" {
		return &api.ServiceInstanceResponse{}, ErrInvalid
	}
	if addRoutingKeyRequest.RoutingKey == "" {
		return &api.ServiceInstanceResponse{}, ErrNotFound
	}

	instance := &routingkeys.ServiceInstance{addRoutingKeyRequest.RoutingKey,
		addRoutingKeyRequest.ServiceUri,
		[]string{addRoutingKeyRequest.Tags}}

	bs.rksrv.Add(instance)

	return &api.ServiceInstanceResponse{instance.Uri,
		instance.Tags[0]}, nil

}

// Delete a routingKey to a service
func (bs basicElsService) RemoveRoutingKey(ctx context.Context, req *api.DeleteRoutingKeyRequest) (*google_protobuf.Empty, error) {
	if req.ServiceUri == "" {
		return &google_protobuf.Empty{}, ErrInvalid
	}
	if req.RoutingKey == "" {
		return &google_protobuf.Empty{}, ErrInvalid
	}

	err := bs.rksrv.Remove(req.ServiceUri, req.RoutingKey)

	return &google_protobuf.Empty{}, err
}

const RoutingKeyTableName = "routingKeys"

// NewBasicService returns a naïve dynamoDb implementation of Service.
func NewBasicService(tableName string, dynamoAddr string, id string, secret string, token string) ElsService {
	rk := routingkeys.New(tableName, dynamoAddr, id, secret, token)

	return basicElsService{rk}
}

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(ElsService) ElsService
