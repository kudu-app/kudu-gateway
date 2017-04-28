package user

import (
	"github.com/rnd/kudu-gateway/client"
	"github.com/rnd/kudu-gateway/web/router"
)

var kuduClient client.KuduServiceClient

// Domain expose user domain implementation.
var Domain domain

// domain is user domain implementation.
type domain struct{}

// PlugRoute registers user domain routes.
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
