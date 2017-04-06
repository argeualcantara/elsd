package grpc

import (
	"fmt"
	"log"
	"net"

	pb "github.com/hpcwp/els-go/api"
	"github.com/hpcwp/els-go/config"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type MyElsServer struct{}

func (s *MyElsServer) GetServiceInstance(ctx context.Context, in *pb.Entity) (*pb.ServiceInstance, error) {
	return &pb.ServiceInstance{ServiceUri: "http://localhost ", Tags: "rw"}, nil
}

// NewS creates a gRPC server instance
func NewServer() *MyElsServer {
	return &MyElsServer{}
}

// Start setups the router instance and starts the server
func (mys *MyElsServer) Start() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterElsServer(s, NewServer())
	reflection.Register(s)

	cfg := config.Load()
	fmt.Printf("els listening on %s:%d\n", cfg.Address, cfg.Port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
