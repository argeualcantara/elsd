package elssrv

import (
	"time"
	"github.com/galo/els-go/pkg/api"
	"github.com/go-kit/kit/log"
	"golang.org/x/net/context"
)

type serviceLoggingMiddleware struct {
	logger log.Logger
	next   ElsService
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


func (mw serviceLoggingMiddleware) GetServiceInstanceByKey(ctx context.Context, routingKey *api.RoutingKeyRequest) (srvIns *api.ServiceInstanceReponse, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetServiceInstance",
			"routingKey", routingKey, "result", srvIns, "error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.GetServiceInstanceByKey(ctx, routingKey)
}


func (mw serviceLoggingMiddleware) AddRoutingKey(ctx context.Context, addRoutingKeyRequest *api.AddRoutingKeyRequest) (srvIns *api.ServiceInstanceReponse, err error)  {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "AddRoutingKey",
			"serviceInstance", addRoutingKeyRequest, "result", srvIns, "error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.AddRoutingKey(ctx, addRoutingKeyRequest)
}