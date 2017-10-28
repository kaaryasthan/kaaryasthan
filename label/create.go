package controller

import (
	"log"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/label/model"
	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// CreateHandler creates label
func (c *Controller) CreateHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusCreated)

	usr := &user.User{ID: userID}
	if err := c.uds.Valid(usr); err != nil {
		log.Println("Couldn't validate user: ", err)
		http.Error(w, "Couldn't validate user: "+usr.ID, http.StatusUnauthorized)
		return
	}

	lbl := new(label.Label)
	if err := jsonapi.UnmarshalPayload(r.Body, lbl); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prj := &project.Project{ID: lbl.ProjectID}
	if err := c.pds.Valid(prj); err != nil {
		log.Println("Couldn't validate project: "+strconv.Itoa(prj.ID), err)
		http.Error(w, "Couldn't find project: "+strconv.Itoa(prj.ID), http.StatusInternalServerError)
		return
	}

	err := c.ds.Create(usr, lbl)
	if err != nil {
		log.Println("Unable to save data: ", err)
		return
	}

	if err := jsonapi.MarshalPayload(w, lbl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
