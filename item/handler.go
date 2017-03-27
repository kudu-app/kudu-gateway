package item

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rnd/kudu/db/item"
	"github.com/unrolled/render"
)

var r = render.New()
var itemReq item.Item

func Index(w http.ResponseWriter, req *http.Request) {
	var err error
	res := make(map[string]interface{})

	err = itemReq.Index(&res)
	if err != nil {
		log.Print(err)
		r.JSON(w, http.StatusInternalServerError,
			map[string]string{"error": "Could not get items"})
		return
	}
	r.JSON(w, http.StatusOK, res)
}

func Post(w http.ResponseWriter, req *http.Request) {
	var err error

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&itemReq); err != nil {
		r.JSON(w, http.StatusBadGateway,
			map[string]string{"error": "Invalid request payload"})
		return
	}
	defer req.Body.Close()

	if itemReq.Goal == "" {
		r.JSON(w, http.StatusBadRequest,
			map[string]string{"error": "Goal cannot be empty"})
		return
	}

	id, err := itemReq.Add()
	if err != nil {
		log.Print(err)
		r.JSON(w, http.StatusInternalServerError,
			map[string]string{"error": "Could not add new item"})
		return
	}
	r.JSON(w, http.StatusCreated, map[string]string{"created": id})
}

func Get(w http.ResponseWriter, req *http.Request) {
	var err error
	var itemRes item.Item

	vars := mux.Vars(req)

	err = itemReq.Get(vars["id"], &itemRes)
	if err != nil {
		log.Print(err)
		r.JSON(w, http.StatusInternalServerError,
			map[string]string{"error": fmt.Sprintf("Could not get speficied item with id: %s", vars["id"])})
		return
	}

	if itemRes.Created.Time().IsZero() {
		r.JSON(w, http.StatusNotFound,
			map[string]string{"error": fmt.Sprintf("Could not find specified item with id: %s", vars["id"])})
		return
	}
	r.JSON(w, http.StatusOK, itemRes)
}
