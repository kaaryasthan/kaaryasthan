package controller

import (
	"log"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	item "github.com/kaaryasthan/kaaryasthan/item/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
	"github.com/pkg/errors"
)

// ListCommentHandler list items
func (c *CommentController) ListCommentHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	usr := &user.User{ID: userID}
	if err := c.uds.Valid(usr); err != nil {
		log.Println("Couldn't validate user: "+usr.ID, err)
		http.Error(w, "Couldn't validate user: "+usr.ID, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)

	vars := mux.Vars(r)
	num := vars["number"]

	var err error
	number, err := strconv.Atoi(num)
	if err != nil {
		log.Printf("Invalid number: %s \n%+v\n", num, errors.WithStack(err))
		http.Error(w, "Invalid number: "+num, http.StatusUnauthorized)
		return
	}

	var objs []*item.Comment
	if objs, err = c.ds.List(number); err != nil {
		log.Println("Couldn't find projects: ", err)
		http.Error(w, "Couldn't find projects: ", http.StatusInternalServerError)
		return
	}

	if err := jsonapi.MarshalPayload(w, objs); err != nil {
		log.Println("Couldn't unmarshal: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
