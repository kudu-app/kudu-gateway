package task

import (
	"github.com/rnd/kudu-gateway/client"
	"github.com/rnd/kudu-gateway/web/router"
)

var kuduClient client.KuduServiceClient

// Domain expose task domain implementation.
var Domain domain

// domain is task domain implementation.
type domain struct{}

// PlugRoute registers tasl domain routes.
func (d domain) PlugRoute(route *router.Router) {
	for _, r := range routes {
		route.R.
			Methods(r.Method).
			Path(r.Path).
			HandlerFunc(r.Handler)
	}
}

// PlugClient attach gRPC client service to user domain.
func (d domain) PlugClient(k client.KuduServiceClient) {
	kuduClient = k
}
