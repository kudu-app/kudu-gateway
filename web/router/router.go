package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Router holds routes information and methods to operate.
type Router struct {
	R *mux.Router
}

// Route contains information for each of the kudu web route.
type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// New create new instance of kudu web router.
func New() *Router {
	router := new(Router)
	router.R = mux.NewRouter().StrictSlash(true)

	return router
}

// Run is a wrapper for http.ListenAndServe.
func (router *Router) Run(addr string) {
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(router.R)

	http.ListenAndServe(addr, n)
}
