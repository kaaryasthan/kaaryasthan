package auth

import (
	"log"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	usrDS := user.NewDatastore(db.DB)
	obj := new(user.User)
	if err := jsonapi.UnmarshalPayload(r.Body, obj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)

	err := usrDS.Create(obj)
	if err != nil {
		log.Println("Unable save data: ", err)
		return
	}
	obj.Password = ""
	if err := jsonapi.MarshalPayload(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
