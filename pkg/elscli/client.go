// Package grpc provides a gRPC client for the add service.
package grpc

import (
	jujuratelimit "github.com/juju/ratelimit"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/sony/gobreaker"
	"google.golang.org/grpc"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/galo/els-go/pkg/elssrv"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/galo/els-go/pkg/api"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"time"

)

func New(conn *grpc.ClientConn, tracer stdopentracing.Tracer, logger log.Logger) elssrv.ElsService {
	// We construct a single ratelimiter middleware, to limit the total outgoing
	// QPS from this client to all methods on the remote instance. We also
	// construct per-endpoint circuitbreaker middlewares to demonstrate how
	// that's done, although they could easily be combined into a single breaker
	// for the entire remote instance, too.

	limiter := ratelimit.NewTokenBucketLimiter(jujuratelimit.NewBucketWithRate(100, 100))

	var getServiceInstanceEndpoint endpoint.Endpoint
	{
		getServiceInstanceEndpoint = grpctransport.NewClient(
			conn,
			"ElsService",
			"GetServiceInstance",
			elssrv.EncodeGRPCGetServiceInstanceResponse,
			elssrv.DecodeGRPGetServiceInstanceRequest,
			api.ServiceInstance{},
			grpctransport.ClientBefore(opentracing.ToGRPCRequest(tracer, logger)),
		).Endpoint()
		getServiceInstanceEndpoint = opentracing.TraceClient(tracer, "GetServiceInstance")(getServiceInstanceEndpoint)
		getServiceInstanceEndpoint = limiter(getServiceInstanceEndpoint)
		getServiceInstanceEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Sum",
			Timeout: 30 * time.Second,
		}))(getServiceInstanceEndpoint)
	}

	return  elssrv.Endpoints {
		GetServiceInstanceEndpoint: getServiceInstanceEndpoint,
	}
}