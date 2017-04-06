package service

import ("log"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC AddServer.
func MakeGRPCServer(endpoints Endpoints, tracer stdopentracing.Tracer, logger log.Logger) pb.ElsServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorLogger(logger),
	}
	return &grpcServer{
		sum: grpctransport.NewServer(
			endpoints.GetSrvInstEndpoint,
			DecodeGRPCSumRequest,
			EncodeGRPCSumResponse,
			append(options, grpctransport.ServerBefore(opentracing.FromGRPCRequest(tracer, "Sum", logger)))...,
		),
		concat: grpctransport.NewServer(
			endpoints.ConcatEndpoint,
			DecodeGRPCConcatRequest,
			EncodeGRPCConcatResponse,
			append(options, grpctransport.ServerBefore(opentracing.FromGRPCRequest(tracer, "Concat", logger)))...,
		),
	}
}