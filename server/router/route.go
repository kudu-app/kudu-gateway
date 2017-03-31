package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

func Init() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	return r
}
