package item

import (
	"net/http"

	"github.com/unrolled/render"
)

var r = render.New()

func Index(w http.ResponseWriter, req *http.Request) {
	r.JSON(w, http.StatusOK, map[string]string{"index": "index"})
}

func Get(w http.ResponseWriter, req *http.Request) {
	r.JSON(w, http.StatusOK, map[string]string{"get": "get"})
}

func Post(w http.ResponseWriter, req *http.Request) {
	r.JSON(w, http.StatusOK, map[string]string{"post": "post"})
}
