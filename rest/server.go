package rest

import (
	"fmt"
	"net/http"

	"github.com/hpcwp/els-go/config"
	dynamodbRoutingKeys "github.com/hpcwp/els-go/dynamodb/routingkeys"
	"github.com/hpcwp/els-go/elserror"

	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
)

const (
	applicationJson = "application/json"
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

	// System handling
	router.PanicHandler = serverPanic

	cfg := config.Load()
	fmt.Printf("els listening on %s:%d\n", cfg.Address, cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Address, cfg.Port), router))
}

func serverPanic(w http.ResponseWriter, r *http.Request, rcv interface{}) {
	err := fmt.Sprintf("%v", rcv)
	log.Error("panic", "err", err)

	w.Header().Set("Content-Type", applicationJson)

	json := fmt.Sprintf("{\"error\":\"%s\",\"code\":%d}", elserror.GeneralError.Message, elserror.GeneralError.Code)
	http.Error(w, json, elserror.GeneralError.StatusCode)
}
