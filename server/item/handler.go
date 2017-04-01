package item

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

var response = render.New()
var item Item

// Index retrieves list of item.
func Index(w http.ResponseWriter, r *http.Request) {
	var err error

	res := make(map[string]interface{})
	err = item.Index(&res)
	if err != nil {
		log.Print(err)
		response.JSON(w, http.StatusInternalServerError,
			map[string]string{"error": "Could not get items"})
		return
	}
	response.JSON(w, http.StatusOK, res)
}

// Post adds new item.
func Post(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Body == nil {
		response.JSON(w, http.StatusBadGateway,
			map[string]string{"error": "Invalid request payload"})
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&item); err != nil {
		response.JSON(w, http.StatusBadGateway,
			map[string]string{"error": "Invalid request payload"})
		return
	}
	defer r.Body.Close()

	if item.Goal == "" {
		response.JSON(w, http.StatusBadRequest,
			map[string]string{"error": "Goal cannot be empty"})
		return
	}

	id, err := item.Add()
	if err != nil {
		log.Print(err)
		response.JSON(w, http.StatusInternalServerError,
			map[string]string{"error": "Could not add new item"})
		return
	}
	response.JSON(w, http.StatusCreated, map[string]string{"created": id})
}

// Get retrieves specified item based on request parameter.
func Get(w http.ResponseWriter, r *http.Request) {
	var err error
	var res Item

	vars := mux.Vars(r)

	err = item.Get(vars["id"], &res)
	if err != nil {
		log.Print(err)
		response.JSON(w, http.StatusInternalServerError,
			map[string]string{"error": fmt.Sprintf("Could not get speficied item with id: %s", vars["id"])})
		return
	}

	if res.Created.Time().IsZero() {
		response.JSON(w, http.StatusNotFound,
			map[string]string{"error": fmt.Sprintf("Could not find specified item with id: %s", vars["id"])})
		return
	}
	response.JSON(w, http.StatusOK, res)
}
