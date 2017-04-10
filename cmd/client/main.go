package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/galo/els-go/pkg/elscli"
	"google.golang.org/grpc"
	"github.com/galo/els-go/pkg/api"
)

func main() {
	// The elscli presumes no service discovery system, and expects users to
	// provide the direct address of an elssvc. This presumption is reflected in
	// the elscli binary and the the client packages: the -transport.addr flags
	// and various client constructors both expect host:port strings. For an
	// example service with a client built on top of a service discovery system,
	// see profilesvc.

	var (
		grpcAddr        = flag.String("grpc.addr", "", "gRPC (HTTP) address of elssvc")
		method          = flag.String("method", "getServiceInstance", "getServiceInstance routingKey")
	)
	flag.Parse()

	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "usage: elscli [flags] <routingKey> \n")
		os.Exit(1)
	}



	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	client :=api.NewElsClient(conn)

	//grpcclient.GetServiceInstanceByKey()

	switch *method {
	case "getServiceInstance":
		routingKey := flag.Args()[0]

		v, err := elscli.GetServiceInstanceByKey(client, routingKey)

		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%d  %d\n", routingKey, v)

	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", method)
		os.Exit(1)
	}

}
