package item

import (
	"net/http"

	"github.com/rnd/kudu/server/router"
)

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
