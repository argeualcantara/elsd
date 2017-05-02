/**
 * (C) Copyright 2012-2016 HP Development Company, L.P.
 * Confidential computer software. Valid license from HP required for possession, use or copying.
 * Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
 * Computer Software Documentation, and Technical Data for Commercial Items are licensed
 * to the U.S. Government under vendor's standard commercial license.
 */
package elssrv

import (
	"github.com/go-kit/kit/metrics"
	"github.com/hpcwp/elsd/pkg/api"
	"golang.org/x/net/context"
)

// ServiceInstrumentingMiddleware returns a service middleware that instruments
// the number of routingKeys accessed over the lifetime of
// the service.
func ServiceInstrumentingMiddleware(ints metrics.Counter) Middleware {
	return func(next ElsService) ElsService {
		return serviceInstrumentingMiddleware{
			ints: ints,
			next: next,
		}
	}
}

type serviceInstrumentingMiddleware struct {
	ints metrics.Counter
	next ElsService
}

func (mw serviceInstrumentingMiddleware) GetServiceInstanceByKey(ctx context.Context, routingKey *api.RoutingKeyRequest) (*api.ServiceInstanceReponse, error) {
	v, err := mw.next.GetServiceInstanceByKey(ctx, routingKey)
	mw.ints.Add(1)
	return v, err
}

func (mw serviceInstrumentingMiddleware) AddRoutingKey(ctx context.Context, addRoutingKeyRequest *api.AddRoutingKeyRequest) (*api.ServiceInstanceReponse, error) {
	v, err := mw.next.AddRoutingKey(ctx, addRoutingKeyRequest)
	mw.ints.Add(1)
	return v, err
}
