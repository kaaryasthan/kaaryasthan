package controller

import (
	"log"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/item/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// CreateCommentHandler creates comment
func (c *CommentController) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusCreated)

	usr := &user.User{ID: userID}
	if err := c.uds.Valid(usr); err != nil {
		log.Println("Couldn't validate user: ", err)
		return
	}

	cmt := new(item.Comment)
	if err := jsonapi.UnmarshalPayload(r.Body, cmt); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	num := vars["number"]

	number, err := strconv.Atoi(num)
	if err != nil {
		log.Println("Invalid number: "+num, err)
		http.Error(w, "Invalid number: "+num, http.StatusUnauthorized)
		return
	}

	itm := &item.Item{ID: number}
	if err := c.ids.Valid(itm); err != nil {
		log.Println("Couldn't validate item: ", err)
		return
	}

	cmt.ItemID = number
	err = c.ds.Create(usr, cmt)
	if err != nil {
		log.Println("Unable to save data: ", err)
		return
	}

	if err := jsonapi.MarshalPayload(w, cmt); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
