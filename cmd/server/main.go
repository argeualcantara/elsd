package main

import (
	"runtime"
	"strings"

	"github.com/dimiro1/banner"
	"os"
	"flag"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/metrics"
	"github.com/galo/els-go/pkg/elssrv"
	"syscall"
	"fmt"
	"os/signal"
	"net/http"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http/pprof"
	"net"
	"google.golang.org/grpc"
	"github.com/galo/els-go/pkg/api"
)

const (
	bannerTxt = `
 ______     __         ______
/\  ___\   /\ \       /\  ___\
\ \  __\   \ \ \____  \ \___  \
 \ \_____\  \ \_____\  \/\_____\
  \/_____/   \/_____/   \/_____/

CWP Entity Locator Service v1.5.0
(C) Copyright 2016-2017 HP Development Company, L.P.

GoVersion: {{ .GoVersion }}
NumCPU: {{ .NumCPU }}
Now: {{ .Now "Mon, 02 Jan 2006 15:04:05 -0700" }}
Debug: '{{ .Env "ELS_DEBUG" }}'
`
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		debugAddr = flag.String("debug.addr", ":8080", "Debug and metrics listen address")
		grpcAddr  = flag.String("grpc.addr", ":8082", "gRPC (HTTP) listen address")
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
	banner.Init(os.Stdout, true, false, strings.NewReader(bannerTxt))

	// Metrics domain.
	var ints metrics.Counter
	{
		// Business level metrics.
		ints = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "addsvc",
			Name:      "integers_summed",
			Help:      "Total count of integers summed via the Sum method.",
		}, []string{})
	}

	// Business domain.
	var service elssrv.ElsService
	{
		service = elssrv.NewBasicService()
		service = elssrv.ServiceLoggingMiddleware(logger)(service)
		service = elssrv.ServiceInstrumentingMiddleware(ints)(service)
	}

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

	// gRPC transport.
	go func() {
		logger := log.With(logger, "transport", "gRPC")

		ln, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errc <- err
			return
		}

		s := grpc.NewServer()
		api.RegisterElsServer(s, service)

		logger.Log("addr", *grpcAddr)
		errc <- s.Serve(ln)
	}()

	// Run!
	logger.Log("exit", <-errc)

}
