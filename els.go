package main

import (
	"os"
	"runtime"
	"os/signal"
	"syscall"
	"strings"

	"github.azc.ext.hp.com/cwp/els-go/rest"
	"github.azc.ext.hp.com/cwp/els-go/config"
	"github.com/dimiro1/banner"
	log "github.com/Sirupsen/logrus"
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
	if (cfg.IsDebug) {
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
