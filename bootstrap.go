package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/knq/firebase"
	"github.com/rnd/kudu/db"
	"github.com/rnd/kudu/item"
	"github.com/rnd/kudu/router"
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

	a.config, err = registerConfig()
	if err != nil {
		log.Fatal(err)
	}

	a.route = registerRoutes()

	err = setupFirebase()
	if err != nil {
		log.Fatal(err)
	}
}

func registerConfig() (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.AddConfigPath(".")
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

func setupFirebase() error {
	var err error
	db.ItemRef, err = firebase.NewDatabaseRef(
		firebase.GoogleServiceAccountCredentialsFile(kudu.config.GetString("firebase.item.cred")),
	)
	return err
}
