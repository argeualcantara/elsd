package elssrv

import (
	"github.com/galo/els-go/pkg/api"
	"golang.org/x/net/context"
	"github.com/go-kit/kit/metrics"
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

