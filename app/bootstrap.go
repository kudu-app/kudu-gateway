package app

import (
	"fmt"
	"log"

	"github.com/knq/envcfg"
	"github.com/rnd/kudu-web/app/router"
)

// Kudu contains information needed to start kudu web application.
type Kudu struct {
	config *envcfg.Envcfg
	router *router.Router
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
	r := router.New()
	// r.RegisterRoutes(home.Routes)
	// r.RegisterRoutes(items.Routes)
	k.router = r

	//TODO: Setup grpc client connection.
}
