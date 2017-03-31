package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route contains information for each of the kudu web route.
type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// Init initialize gorilla mux with strict slash enabled.
func Init() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	return r
}
