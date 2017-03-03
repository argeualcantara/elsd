package main

import (
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/dimiro1/banner"
	"github.com/hpcwp/els-go/config"
	"github.com/hpcwp/els-go/rest"
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
	// Shows fancy banner
	banner.Init(os.Stdout, true, false, strings.NewReader(bannerTxt))

	// Load up configs and setup logging
	cfg := config.Load()
	if cfg.IsDebug {
		log.SetLevel(log.DebugLevel)
	}

	log.SetOutput(os.Stdout)

	log.Info("ELS is starting...")
	server := rest.New()
	go server.Start()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Warn("received system signal", "signal", <-ch)
}
