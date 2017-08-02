/**
 * (C) Copyright 2012-2016 HP Development Company, L.P.
 * Confidential computer software. Valid license from HP required for possession, use or copying.
 * Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
 * Computer Software Documentation, and Technical Data for Commercial Items are licensed
 * to the U.S. Government under vendor's standard commercial license.
 */

package elssrv

import (
	"github.com/hpcwp/elsd/pkg/dynamodb/routingkeys"
)

type ElsService struct {
	rk *routingkeys.Service
}

type ServiceInstance struct {
	Url      string `json:"url"`
	Metadata string `json:"metadata"`
}

func (s ElsService) AddKey(key string, srv ServiceInstance) error {
	if srv.Url == "" {
		return ErrInvalid
	}

	if key == "" {
		return ErrInvalid
	}

	instance := &routingkeys.ServiceInstance{key,
		srv.Url,
		[]string{srv.Metadata}}

	s.rk.Add(instance)

	return nil
}

func (s ElsService) RemoveService(key string, srvUri string) error {
	if srvUri == "" {
		return ErrInvalid
	}
	if key == "" {
		return ErrInvalid
	}

	err := s.rk.Remove(srvUri, key)

	return err
}

func (s ElsService) GetService(key string) (*ServiceInstance, error) {

	if key == "" {
		return nil, ErrInvalid
	}

	serviceInstance := s.rk.Get(key)

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

	srvInstance := ServiceInstance{serviceUrl, "rw"}

	return &srvInstance, nil
}

func (s ElsService) ListServices(key string) ([]*ServiceInstance, error) {
	if key == "" {
		return make([]*ServiceInstance, 0), ErrInvalid
	}

	entities := s.rk.Get(key)

	if entities == nil {
		return make([]*ServiceInstance, 0), ErrNotFound
	}
	if len(entities.ServiceInstances) == 0 {
		return make([]*ServiceInstance, 0), ErrNotFound
	}

	listResp := make([]*ServiceInstance, 0)

	for i := range entities.ServiceInstances {
		srvInstance := ServiceInstance{entities.ServiceInstances[i].Uri, "rw"}
		listResp = append(listResp, &srvInstance)

	}

	return listResp, nil
}

// NewBasicService returns a naïve dynamoDb implementation of Service.
func NewService(tableName string, dynamoAddr string, id string, secret string, token string) ElsService {
	rk := routingkeys.New(tableName, dynamoAddr, id, secret, token)

	return ElsService{rk}
}
