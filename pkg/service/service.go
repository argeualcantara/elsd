package service

// This file contains the Service definition, and a basic service
// implementation. It also includes service middlewares.

import (
	"errors"
)

// Service describes a service that adds things together.
type ElsService interface {
	GetServiceInstance(routingKey string) (ServiceInstance, error)
}


type ServiceInstance struct {
	Url string `json:"url"`
	Metadata string `json:"metadata"`
}


type elsService struct {}

// The implementation of the service
func (elsService) GetServiceInstance( routingKey string) (ServiceInstance, error) {
	if routingKey =="" {
		return nil, ErrInvalid
	}
	srvInstance := ServiceInstance{"http://localhost","rw"}
	return srvInstance, nil
}



// ErrEmpty is returned when input is invalid
var ErrInvalid = errors.New("Invalid routing key")

