package web

import (
	"fmt"
	"log"

	"github.com/knq/envcfg"

	"github.com/rnd/kudu-gateway/client"
	"github.com/rnd/kudu-gateway/web/domain"
	"github.com/rnd/kudu-gateway/web/domain/task"
	"github.com/rnd/kudu-gateway/web/domain/user"
	"github.com/rnd/kudu-gateway/web/router"
)

// Kudu contains information needed to start kudu web application.
type Kudu struct {
	config *envcfg.Envcfg
	router *router.Router
	client client.KuduServiceClient
}

// Run starts kudu web application server.
func (k *Kudu) Run() {
	k.bootstrap()

	addr := fmt.Sprintf(":%s", k.config.GetString("app.port"))
	log.Printf("Application start on: %s \n", addr)
	k.router.Run(addr)
}

// bootstrap prepare and do preprocess works before actually running kudu application.
func (k *Kudu) bootstrap() {
	var err error

	// setup app config.
	k.config, err = envcfg.New()
	if err != nil {
		log.Fatal(err)
	}

	// setup app routes.
	k.router = router.New()

	// grpc client connections.
	k.client, err = client.New()
	if err != nil {
		log.Fatal(err)
	}

	// register domain services.
	k.plug(
		user.Domain,
		task.Domain,
	)
}

// plug registers all kudu domains to kudu web api gateway.
func (k *Kudu) plug(domains ...domain.Domain) {
	for _, domain := range domains {
		domain.PlugRoute(k.router)
		domain.PlugClient(k.client)
	}
}
