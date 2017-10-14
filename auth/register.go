package auth

import (
	"log"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	obj := new(user.User)
	if err := jsonapi.UnmarshalPayload(r.Body, obj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)

	err := obj.Create()
	if err != nil {
		log.Println("Unable save data: ", err)
		return
	}
	obj.Password = ""
	if err := jsonapi.MarshalPayload(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
