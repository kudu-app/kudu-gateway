package domain

import (
	"github.com/rnd/kudu-gateway/client"
	"github.com/rnd/kudu-gateway/web/router"
)

// Domain is standard interface for all kudu domains.
type Domain interface {
	// PlugRoute registers domain routes.
	PlugRoute(r *router.Router)

	// PlugClient attach gRPC client service to each domain.
	PlugClient(k client.KuduServiceClient)
}
