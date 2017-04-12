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

func (mw serviceLoggingMiddleware) GetServiceInstanceByKey(ctx context.Context, routingKey *api.RoutingKey) (srvIns *api.ServiceInstance, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetServiceInstance",
			"routingKey", routingKey, "result", srvIns, "error", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return mw.next.GetServiceInstanceByKey(ctx, routingKey)
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