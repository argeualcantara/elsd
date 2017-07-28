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
	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/hpcwp/elsd/pkg/api"
	"golang.org/x/net/context"
)

type serviceInstrumentingMiddleware struct {
	keys    metrics.Gauge
	queries metrics.Counter
	next    GRPCServer
}

// ServiceInstrumentingMiddleware returns a service middleware that instruments
// the number of routingKeys accessed over the lifetime of
// the service.
func ServiceInstrumentingMiddleware(keys metrics.Gauge, queries metrics.Counter) Middleware {
	return func(next GRPCServer) GRPCServer {
		return serviceInstrumentingMiddleware{
			keys:    keys,
			queries: queries,
			next:    next,
		}
	}
}

func (mw serviceInstrumentingMiddleware) GetServiceInstanceByKey(ctx context.Context, routingKey *api.RoutingKeyRequest) (*api.ServiceInstanceResponse, error) {
	v, err := mw.next.GetServiceInstanceByKey(ctx, routingKey)
	mw.queries.Add(1)
	return v, err
}

func (mw serviceInstrumentingMiddleware) ListServiceInstances(ctx context.Context, routingKey *api.RoutingKeyRequest) (*api.ServiceInstanceListResponse, error) {
	v, err := mw.next.ListServiceInstances(ctx, routingKey)
	mw.queries.Add(1)
	return v, err
}

func (mw serviceInstrumentingMiddleware) AddRoutingKey(ctx context.Context, addRoutingKeyRequest *api.AddRoutingKeyRequest) (*api.ServiceInstanceResponse, error) {
	v, err := mw.next.AddRoutingKey(ctx, addRoutingKeyRequest)
	mw.queries.Add(1)
	mw.keys.Add(1)
	return v, err
}

func (mw serviceInstrumentingMiddleware) RemoveRoutingKey(ctx context.Context, req *api.DeleteRoutingKeyRequest) (empty *google_protobuf.Empty, err error) {
	v, err := mw.next.RemoveRoutingKey(ctx, req)
	mw.queries.Add(1)
	// Since delete is an idempotent operation it is possible that the gauge becomes negative
	if err != nil {
		mw.keys.Add(-1)
	}
	return v, err

}
