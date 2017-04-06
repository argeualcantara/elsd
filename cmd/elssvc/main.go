package main

import (
	"runtime"
	"strings"

	"github.com/dimiro1/banner"
	"os"
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
}


