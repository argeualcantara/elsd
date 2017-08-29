/**
 * (C) Copyright 2012-2016 HP Development Company, L.P.
 * Confidential computer software. Valid license from HP required for possession, use or copying.
 * Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
 * Computer Software Documentation, and Technical Data for Commercial Items are licensed
 * to the U.S. Government under vendor's standard commercial license.
 */
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hpcwp/elsd/pkg/api"
	"github.com/hpcwp/elsd/pkg/elscli"
	"google.golang.org/grpc"
	"os"
	"time"
)

const defaulGrpcAddres string = "localhost:8082"

func main() {
	// The elscli presumes no service discovery system, and expects users to
	// provide the direct address of an elssvc.

	var (
		grpcAddr = flag.String("grpc.addr", "localhost:8082", "gRPC (HTTP) address of elssvc")
		method   = flag.String("method", "Check", "Get routingKey, Add routingKey uri tags")
	)
	flag.Parse()

	if method == nil {
		fmt.Fprintf(os.Stderr, "usage: elscli -grpc.addr <address> -method {Get|Add} <arguments> \n")
		os.Exit(1)
	}

	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	client := api.NewElsClient(conn)
	healthClient := api.NewHealthClient(conn)

	switch *method {
	case "Check":
		check(healthClient)

	case "Get":
		get(client)

	case "List":
		list(client)

	case "Add":
		add(client)

	case "Remove":
		remove(client)

	default:
		fmt.Fprintf(os.Stderr, "error: invalid method %q\n", method)
		os.Exit(1)
	}

}
func remove(client api.ElsClient) {
	if len(flag.Args()) != 2 {
		fmt.Fprintf(os.Stderr, "usage: elscli -grpc.addr <address> -method Remove routing-key uri \n")
		os.Exit(1)
	}
	routingKey := flag.Args()[0]
	uri := flag.Args()[1]
	v, err := elscli.RemoveServiceInstance(client, routingKey, uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%d  %d\n", routingKey, v)
}
func add(client api.ElsClient) {
	if len(flag.Args()) != 3 {
		fmt.Fprintf(os.Stderr, "usage: elscli -grpc.addr <address> -method Add routing-key uri tags\n")
		os.Exit(1)
	}
	routingKey := flag.Args()[0]
	uri := flag.Args()[1]
	tags := flag.Args()[2]
	v, err := elscli.AddServiceInstance(client, routingKey, uri, []string{tags})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%d  %d\n", routingKey, v)
}
func list(client api.ElsClient) {
	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "usage: elscli -grpc.addr <address> -method Get routing-key \n")
		os.Exit(1)
	}
	routingKey := flag.Args()[0]
	v, err := elscli.ListServiceInstances(client, routingKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%d  %d\n", routingKey, v)
}
func get(client api.ElsClient) {
	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "usage: elscli -grpc.addr <address> -method Get routing-key \n")
		os.Exit(1)
	}
	routingKey := flag.Args()[0]
	v, err := elscli.GetServiceInstanceByKey(client, routingKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%d  %d\n", routingKey, v)
}

func check(healthClient api.HealthClient) {
	_, err := healthClient.Check(context.Background(), &api.HealthCheckRequest{"els"})
	if err != nil {
		fmt.Fprintf(os.Stderr, "elsd service is not responding \n")
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "elsd is ready \n")
}
