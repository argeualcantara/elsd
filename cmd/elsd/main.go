/**
 * (C) Copyright 2012-2016 HP Development Company, L.P.
 * Confidential computer software. Valid license from HP required for possession, use or copying.
 * Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
 * Computer Software Documentation, and Technical Data for Commercial Items are licensed
 * to the U.S. Government under vendor's standard commercial license.
 */
package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/dimiro1/banner"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/hpcwp/elsd/pkg/api"
	"github.com/hpcwp/elsd/pkg/elssrv"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const (
	bannerTxt = `
{{.AnsiColor.Red}}EEEEEEEEEEEEEEEEEEEEEELLLLLLLLLLL                SSSSSSSSSSSSSSS             d::::::d
E::::::::::::::::::::EL:::::::::L              SS:::::::::::::::S            d::::::d
E::::::::::::::::::::EL:::::::::L             S:::::SSSSSS::::::S            d::::::d
{{.AnsiColor.BrightBlue}}EE::::::EEEEEEEEE::::ELL:::::::LL             S:::::S     SSSSSSS            d:::::d
  E:::::E       EEEEEE  L:::::L               S:::::S                ddddddddd:::::d
  E:::::E               L:::::L               S:::::S              dd::::::::::::::d
  E::::::EEEEEEEEEE     L:::::L                S::::SSSS          d::::::::::::::::d
  E:::::::::::::::E     L:::::L                 SS::::::SSSSS    d:::::::ddddd:::::d
  E:::::::::::::::E     L:::::L                   SSS::::::::SS  d::::::d    d:::::d
  E::::::EEEEEEEEEE     L:::::L                      SSSSSS::::S d:::::d     d:::::d
  E:::::E               L:::::L                           S:::::Sd:::::d     d:::::d
  E:::::E       EEEEEE  L:::::L         LLLLLL            S:::::Sd:::::d     d:::::d
{{ .AnsiColor.Default }}EE::::::EEEEEEEE:::::ELL:::::::LLLLLLLLL:::::LSSSSSSS     S:::::Sd::::::ddddd::::::dd
E::::::::::::::::::::EL::::::::::::::::::::::LS::::::SSSSSS:::::S d:::::::::::::::::d
E::::::::::::::::::::EL::::::::::::::::::::::LS:::::::::::::::SS   d:::::::::ddd::::d
EEEEEEEEEEEEEEEEEEEEEELLLLLLLLLLLLLLLLLLLLLLLL SSSSSSSSSSSSSSS      ddddddddd   ddddd

CWP Entity Locator Service v2.0.0
(C) Copyright 2016-2017 HP Development Company, L.P.

GoVersion: {{ .GoVersion }}
GOOS: {{ .GOOS }}
GOARCH: {{ .GOARCH }}
NumCPU: {{ .NumCPU }}
GOPATH: {{ .GOPATH }}
GOROOT: {{ .GOROOT }}
Compiler: {{ .Compiler }}
ENV: {{ .Env "GOPATH" }}
`
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		debugAddr    = flag.String("debug.addr", ":8080", "Debug and metrics listen address")
		grpcAddr     = flag.String("grpc.addr", ":8082", "gRPC (HTTP) listen address")
		dynamoDbAddr = flag.String("dynamodb.addr", "http://localhost:8000", "DynamoDb (HTTP/HTTPS) address")
		region       = flag.String("aws.region", "us-west-2", "AWS dynamo region")
		id           = flag.String("aws.id", "123", "AWS id")
		secret       = flag.String("aws.secret", "123", "AWS secret")
		token        = flag.String("aws.token", "", "AWS token")
	)

	flag.Parse()

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Shows fancy banner
	banner.Init(os.Stdout, true, true, strings.NewReader(bannerTxt))

	logger.Log("debugAddr", debugAddr, "grpcAddr", grpcAddr, "dynamodbAddr", dynamoDbAddr)

	// Queries domain.
	var queries metrics.Counter
	{
		// Business level metrics.
		queries = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "elds",
			Name:      "queries",
			Help:      "Total queries.",
		}, []string{})
	}

	// Metrics domain.
	var keys metrics.Gauge
	{
		// Business level metrics.
		keys = prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: "elds",
			Name:      "keys",
			Help:      "Keys stored.",
		}, []string{})
	}

	// Business domain.
	var service elssrv.GRPCServer
	{
		service = elssrv.NewBasicService(elssrv.RoutingKeyTableName, *dynamoDbAddr, *region, *id, *secret, *token)
		service = elssrv.ServiceLoggingMiddleware(logger)(service)
		service = elssrv.ServiceInstrumentingMiddleware(keys, queries)(service)
	}

	var healthService elssrv.HealthGRPCServer

	// Mechanical domain.
	errc := make(chan error)

	// Interrupt handler.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// Debug listener.
	go func() {
		logger := log.With(logger, "transport", "debug")

		m := http.NewServeMux()
		m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		m.Handle("/metrics", promhttp.Handler())

		logger.Log("addr", *debugAddr)
		errc <- http.ListenAndServe(*debugAddr, m)
	}()

	// gRPC elsd transport.
	go func() {
		logger := log.With(logger, "transport", "gRPC")

		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}

		s := grpc.NewServer()
		api.RegisterElsServer(s, service)
		api.RegisterHealthServer(s, healthService)

		logger.Log("addr", grpcAddr)
		errc <- s.Serve(ln)
	}()

	// Run!
	logger.Log("exit", <-errc)

}
