package main

import (
	"flag"
	"os"
	"fmt"
	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/galo/els-go/pkg/elssrv"
	grpcclient "github.com/galo/els-go/pkg/elscli"
	"github.com/go-kit/kit/log"
	"strings"
	"google.golang.org/grpc"
	"time"
)

func main() {
	// The addcli presumes no service discovery system, and expects users to
	// provide the direct address of an addsvc. This presumption is reflected in
	// the addcli binary and the the client packages: the -transport.addr flags
	// and various client constructors both expect host:port strings. For an
	// example service with a client built on top of a service discovery system,
	// see profilesvc.

	var (
		grpcAddr        = flag.String("grpc.addr", "", "gRPC (HTTP) address of addsvc")
		zipkinAddr      = flag.String("zipkin.addr", "", "Enable Zipkin tracing via a Zipkin HTTP Collector endpoint")
		zipkinKafkaAddr = flag.String("zipkin.kafka.addr", "", "Enable Zipkin tracing via a Kafka server host:port")
		method          = flag.String("method", "getServiceInstance", "getServiceInstance routingKey")
	)
	flag.Parse()

	if len(flag.Args()) != 2 {
		fmt.Fprintf(os.Stderr, "usage: elscli [flags] <a> <b>\n")
		os.Exit(1)
	}

	var tracer stdopentracing.Tracer
	{
		if *zipkinAddr != "" {
			// endpoint typically looks like: http://zipkinhost:9411/api/v1/spans
			collector, err := zipkin.NewHTTPCollector(*zipkinAddr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			defer collector.Close()

			tracer, err = zipkin.NewTracer(
				zipkin.NewRecorder(collector, false, "0.0.0.0:0", "addcli"),
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		} else if *zipkinKafkaAddr != "" {
			collector, err := zipkin.NewKafkaCollector(
				strings.Split(*zipkinKafkaAddr, ","),
				zipkin.KafkaLogger(log.NewNopLogger()),
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			defer collector.Close()

			tracer, err = zipkin.NewTracer(
				zipkin.NewRecorder(collector, false, "0.0.0.0:0", "addcli"),
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
		} else {
			tracer = stdopentracing.GlobalTracer() // no-op
		}
	}

	var (
		service elssrv.ElsService
		err     error
	)

	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(time.Second))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	service = grpcclient.New(conn, tracer, log.NewNopLogger())

	switch *method {
	case "getServiceInstance":
		routingKey := flag.Args()[0]

		v, err := service.GetServiceInstance(routingKey)
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
