package main

import (
	"runtime"

	"github.azc.ext.hp.com/cwp/els-go/rest"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	address = kingpin.Flag("address", "TCP Address to listen at").Default("localhost").String()
	port    = kingpin.Flag("port", "TCP Port").Default("8080").Int32()
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	kingpin.Parse()

	server := rest.New(*address, *port)
	server.Start()
}
