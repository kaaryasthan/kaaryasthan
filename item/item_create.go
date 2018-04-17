package controller

import (
	"log"
	"net/http"
	"strconv"

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

	usr := &user.User{ID: userID}
	if err := c.uds.Valid(usr); err != nil {
		log.Println("Couldn't validate user: ", err)
		return
	}

	itm := new(item.Item)
	if err := jsonapi.UnmarshalPayload(r.Body, itm); err != nil {
		log.Printf("Couldn't unmarshall item: %#v   %s", itm, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v, err := strconv.Atoi(itm.ProjectID)
	if err != nil {
		log.Println("Unable to convert data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	prj := &project.Project{ID: v}
	if err := c.pds.Valid(prj); err != nil {
		log.Println("Couldn't validate project: ", prj, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.ds.Create(usr, itm)
	if err != nil {
		log.Println("Unable to save data: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := jsonapi.MarshalPayload(w, itm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
