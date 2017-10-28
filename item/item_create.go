package controller

import (
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/item/model"
	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// CreateItemHandler creates item
func (c *ItemController) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusCreated)

	usr := &user.User{ID: userID}
	if err := c.uds.Valid(usr); err != nil {
		log.Println("Couldn't validate user: ", err)
		return
	}

	itm := new(item.Item)
	if err := jsonapi.UnmarshalPayload(r.Body, itm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prj := &project.Project{ID: itm.ProjectID}
	if err := c.pds.Valid(prj); err != nil {
		log.Println("Couldn't validate project: ", err)
		return
	}

	err := c.ds.Create(usr, itm)
	if err != nil {
		log.Println("Unable to save data: ", err)
		return
	}

	if err := jsonapi.MarshalPayload(w, itm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
