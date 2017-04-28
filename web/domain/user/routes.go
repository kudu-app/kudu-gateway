package user

import (
	"net/http"

	"github.com/rnd/kudu-gateway/web/router"
)

// routes list all user domain routes.
var routes = []*router.Route{
	{
		Method:  http.MethodPost,
		Path:    "/user/register",
		Handler: Register,
	},
	{
		Method:  http.MethodPost,
		Path:    "/user/login",
		Handler: Login,
	},
}
