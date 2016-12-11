package rest

import (
	"fmt"
	"log"
	"net/http"

	"github.azc.ext.hp.com/cwp/els-go/config"
	dynamodbRoutingKeys "github.azc.ext.hp.com/cwp/els-go/dynamodb/routingkeys"

	"github.com/julienschmidt/httprouter"
)

// Server provides a type to handle the REST server
type Server struct {
}

// New creates a REST server instance
func New() *Server {
	return &Server{}
}

// Start setups the router instance and starts the server
func (s *Server) Start() {
	router := httprouter.New()
	routingKeysSvc := dynamodbRoutingKeys.New("RoutingKeys")

	// RoutingKeys
	router.GET("/api/v1/routingkeys/:id", routingKeysSvc.RoutingKeysGet)

	cfg := config.Load()
	fmt.Printf("els listening on %s:%d\n", cfg.Address, cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), router))
}
