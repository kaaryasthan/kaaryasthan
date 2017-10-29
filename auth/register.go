package auth

import (
	"log"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// RegisterHandler register user
func (c *Controller) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	usr := new(user.User)
	if err := jsonapi.UnmarshalPayload(r.Body, usr); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)

	err := c.uds.Create(usr)
	if err != nil {
		log.Println("Unable save data: ", err)
		return
	}
	usr.Password = ""
	if err := jsonapi.MarshalPayload(w, usr); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
