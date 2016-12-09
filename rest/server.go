package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.azc.ext.hp.com/cwp/els-go/dynamodb/routingkeys"

	"github.com/julienschmidt/httprouter"
)

// Server provides a type to handle the REST server
type Server struct {
	address string
	port    int32
}

// New creates a REST server instance
func New(address string, port int32) *Server {
	return &Server{
		address: address,
		port:    port,
	}
}

// Start setups the router instance and starts the server
func (s *Server) Start() {
	router := httprouter.New()
	routingKeysSvc := routingkeys.New("RoutingKeys")

	// RoutingKeys
	router.GET("/els/api/v1/routingkeys/:id", routingKeysSvc.RoutingKeysGet)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", s.address, s.port), router))
}
