package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rnd/kudu/server/db"
	"github.com/rnd/kudu/server/item"
	"github.com/rnd/kudu/server/router"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
)

// app contains information needed to start kudu application.
type app struct {
	config *viper.Viper
	route  http.Handler
}

// run starts kudu web server on specified port.
func (a *app) run() {
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(a.route)

	addr := fmt.Sprintf(":%s", a.config.GetString("port"))
	n.Run(addr)
}

// bootstrap prepare and do preprocess works before actually running kudu application.
func (a *app) bootstrap() {
	var err error

	a.config, err = registerConfig("$GOPATH/src/github.com/rnd/kudu/server")
	if err != nil {
		log.Fatal(err)
	}

	a.route = registerRoutes()

	err = setupDatabase(a.config.GetString("google.creds.data"))
	if err != nil {
		log.Fatal(err)
	}
}

// registerConfig registers kudu configuration file.
func registerConfig(path string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(path)
	return v, v.ReadInConfig()
}

// registerRoutes registers kudu web routes.
func registerRoutes() http.Handler {
	r := router.Init()

	// register item routes.
	for _, v := range item.Routes() {
		r.
			Methods(v.Method).
			Path(v.Path).
			HandlerFunc(v.Handler)
	}

	return r
}

// setupDatabase setup firebase database.
func setupDatabase(path string) error {
	return db.Setup(path)
}
