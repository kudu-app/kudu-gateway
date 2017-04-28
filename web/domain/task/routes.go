package task

import (
	"net/http"

	"github.com/rnd/kudu-gateway/web/router"
)

// routes list all task domain routes.
var routes = []*router.Route{
	{
		Method:  http.MethodGet,
		Path:    "/task/{period}",
		Handler: List,
	},
	{
		Method:  http.MethodPost,
		Path:    "/task",
		Handler: Add,
	},
	{
		Method:  http.MethodDelete,
		Path:    "/task/{id}",
		Handler: Remove,
	},
}
