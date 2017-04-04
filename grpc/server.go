package grpc


import (
	"log"
	"net"

	pb "github.com/hpcwp/els-go/api"
	"google.golang.org/grpc/reflection"
	"github.com/hpcwp/els-go/config"
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/grpc"
	"context"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}



// New creates a gRPC server instance
func New() *server {
	return &server{}
}

// Start setups the router instance and starts the server
func (s *server) Start() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	routingKeysSvc := dynamodbRoutingKeys.New("RoutingKeys")

	// RoutingKeys
	router.GET("/api/v1/routingkeys/:id", routingKeysSvc.RoutingKeysGet)

	cfg := config.Load()
	fmt.Printf("els listening on %s:%d\n", cfg.Address, cfg.Port)

}


// SayHello implements helloworld.GreeterServer
func (s *server) GeServiceINstance(ctx context.Context, in *pb.Entity) (*pb.ServiceInstance, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}