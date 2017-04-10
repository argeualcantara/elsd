package elssrv

// This file contains the Service definition, and a basic service
// implementation. It also includes service middlewares.

import (
	"errors"
	"github.com/go-kit/kit/log"
	"time"
	"github.com/go-kit/kit/metrics"
	"github.com/galo/els-go/pkg/api"
	"golang.org/x/net/context"
)

// Service describes a service that adds things together.
type ElsService interface {

	GetServiceInstanceByKey(ctx context.Context, routingKey *api.RoutingKey) (*api.ServiceInstance, error)
}

type ServiceInstance struct {
	Url      string `json:"url"`
	Metadata string `json:"metadata"`
}

type basicElsService struct{}

// The implementation of the service
func (basicElsService) GetServiceInstanceByKey(ctx context.Context, routingKey *api.RoutingKey) (*api.ServiceInstance, error) {
	if routingKey.Id == "" {
		return &api.ServiceInstance{}, ErrInvalid
	}
	//TODO: returning a fake instance
	srvInstance := api.ServiceInstance{"http://localhost", "rw"}
	return &srvInstance, nil
}

// NewBasicService returns a naïve, stateless implementation of Service.
func NewBasicService() ElsService {
	return basicElsService{}
}

// ServiceLoggingMiddleware returns a service middleware that logs the
// parameters and result of each method invocation.
func ServiceLoggingMiddleware(logger log.Logger) Middleware {
	return func(next ElsService) ElsService {
		return serviceLoggingMiddleware{
			logger: logger,
			next:   next,
		}
	}
}

type serviceLoggingMiddleware struct {
	logger log.Logger
	next   ElsService
}

func (mw serviceLoggingMiddleware) GetServiceInstanceByKey(ctx context.Context, routingKey *api.RoutingKey) (srvIns *api.ServiceInstance, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetServiceInstance",
			"routingKey", routingKey, "result", srvIns, "error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.GetServiceInstanceByKey(ctx,routingKey)
}

// ServiceInstrumentingMiddleware returns a service middleware that instruments
// the number of routingKeys accesed over the lifetime of
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

func (mw serviceInstrumentingMiddleware) GetServiceInstanceByKey(ctx context.Context, routingKey *api.RoutingKey) (*api.ServiceInstance, error) {
	v, err := mw.GetServiceInstanceByKey(ctx, routingKey)
	mw.ints.Add(1)
	return v, err
}

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(ElsService) ElsService

// ErrEmpty is returned when input is invalid
var ErrInvalid = errors.New("Invalid routing key")
