package routingkeys

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// RoutingKeysGet returns a Routing Key
func (s *Service) RoutingKeysGet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	enc := json.NewEncoder(w)
	enc.Encode(s.Get(p.ByName("id")))
}
