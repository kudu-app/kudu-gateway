package item

import (
	"net/http"

	"github.com/rnd/kudu/server/router"
)

// Routes return all item routes.
func Routes() []*router.Route {
	return []*router.Route{
		{
			Method:  http.MethodGet,
			Path:    "/items",
			Handler: Index,
		},
		{
			Method:  http.MethodPost,
			Path:    "/items",
			Handler: Post,
		},
		{
			Method:  http.MethodGet,
			Path:    "/items/{id}",
			Handler: Get,
		},
	}
}
