/**
 * (C) Copyright 2012-2016 HP Development Company, L.P.
 * Confidential computer software. Valid license from HP required for possession, use or copying.
 * Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
 * Computer Software Documentation, and Technical Data for Commercial Items are licensed
 * to the U.S. Government under vendor's standard commercial license.
 */
package elssrv

import (
	"github.com/go-kit/kit/log"
	"github.com/hpcwp/elsd/pkg/api"
	"golang.org/x/net/context"
	"time"
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