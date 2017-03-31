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

type app struct {
	config *viper.Viper
	route  http.Handler
}

func (a *app) run() {
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(a.route)

	addr := fmt.Sprintf(":%s", a.config.GetString("port"))
	n.Run(addr)
}

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

func registerConfig(path string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(path)
	return v, v.ReadInConfig()
}

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

func setupDatabase(path string) error {
	return db.Setup(path)
}
