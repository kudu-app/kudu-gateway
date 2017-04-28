package user

import (
	"encoding/json"
	"net/http"

	"github.com/unrolled/render"

	pcreds "github.com/rnd/kudu/golang/protogen/type/creds"
	userpb "github.com/rnd/kudu/golang/protogen/user"
)

var response = render.New()

// Register creates new user record.
func Register(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Body == nil {
		response.JSON(w, http.StatusBadGateway,
			map[string]string{"error": "Invalid request payload"})
		return
	}

	var params struct {
		Email     string `json:"email"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		response.JSON(w, http.StatusBadGateway,
			map[string]string{"error": "Invalid request payload"})
		return
	}
	defer r.Body.Close()

	req := &userpb.RegisterRequest{
		Credential: &pcreds.Credential{
			Email:    params.Email,
			Username: params.Username,
			Password: params.Password,
		},
		User: &userpb.User{
			FirstName: params.FirstName,
			LastName:  params.LastName,
		},
	}
	res, err := kuduClient.Register(r.Context(), req)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError,
			map[string]string{"error": err.Error()})
		return
	}
	response.JSON(w, http.StatusCreated, map[string]string{
		"status": res.Status.String(),
	})
}

// Login attempt loggin in user to kudu system.
func Login(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Body == nil {
		response.JSON(w, http.StatusBadGateway,
			map[string]string{"error": "Invalid request payload"})
		return
	}

	var params struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		response.JSON(w, http.StatusBadGateway,
			map[string]string{"error": "Invalid request payload"})
		return
	}
	defer r.Body.Close()

	req := &userpb.LoginRequest{
		Credential: &pcreds.Credential{
			Email:    params.Email,
			Password: params.Password,
		},
	}

	res, err := kuduClient.Login(r.Context(), req)
	if err != nil {
		response.JSON(w, http.StatusBadGateway,
			map[string]string{"error": err.Error()})
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"status": res.Status.String(),
		"token":  res.Token,
	})
}
