package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Router holds routes information and methods to operate.
type Router struct {
	r      *mux.Router
	routes []Route
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
	router.r = mux.NewRouter().StrictSlash(true)

	return router
}

// RegisterRoutes registers kudu web routes.
func (router *Router) RegisterRoutes(routeGroups ...[]*Route) {
	for _, routes := range routeGroups {
		for _, route := range routes {
			router.r.
				Methods(route.Method).
				Path(route.Path).
				HandlerFunc(route.Handler)
		}
	}
}

// Run is a wrapper for http.ListenAndServe.
func (router *Router) Run(addr string) {
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(router.r)

	http.ListenAndServe(addr, n)
}
