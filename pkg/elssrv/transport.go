package elssrv

import (stdopentracing "github.com/opentracing/opentracing-go"
	oldcontext "golang.org/x/net/context"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/galo/els-go/pkg/api"
	"context"
	"github.com/go-kit/kit/log"
)



// MakeGRPCServer makes a set of endpoints available as a gRPC AddServer.
func MakeGRPCServer(endpoints Endpoints, tracer stdopentracing.Tracer, logger log.Logger) api.ElsServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		getServiceInstance: grpctransport.NewServer(
			endpoints.GetSrvInstEndpoint,
			DecodeGRPGetServiceInstanceRequest,
			EncodeGRPCGetServiceInstanceResponse,
			append(options, grpctransport.ServerBefore(opentracing.FromGRPCRequest(tracer, "Sum", logger)))...,
		),
	}
}

type grpcServer struct {
	getServiceInstance    grpctransport.Handler

}

func (s *grpcServer) GetServiceInstance(ctx oldcontext.Context, req *api.Entity) (*api.ServiceInstance, error) {
	_, rep, err := s.getServiceInstance.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*api.ServiceInstance), nil
}


// DecodeGRPGetServiceInstanceRequest is a transport/grpc.EncodeRequestFunc that converts a
// grpc Entity request into a user-domain DecodeGRPGetServiceInstanceRequest
func DecodeGRPGetServiceInstanceRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(api.Entity)

	return getServiceInstanceRequest{req.GetId()}, nil
}

func EncodeGRPCGetServiceInstanceResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(getServiceInstanceResponse)
	return &api.ServiceInstance{resp.ServiceInstance.Url, resp.ServiceInstance.Metadata}, nil
}