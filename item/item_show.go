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
	"github.com/pkg/errors"
)

// ShowItemHandler shows item
func (c *ItemController) ShowItemHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)

	usr := &user.User{ID: userID}
	if err := c.uds.Valid(usr); err != nil {
		log.Printf("Couldn't validate user: %s \n%+v\n", usr.ID, err)
		http.Error(w, "Couldn't validate user: "+usr.ID, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	num := vars["number"]

	number, err := strconv.Atoi(num)
	if err != nil {
		log.Printf("Invalid number: %s \n%+v\n", num, errors.WithStack(err))
		http.Error(w, "Invalid number: "+num, http.StatusUnauthorized)
		return
	}

	itm := &item.Item{Number: num}
	if err := c.ds.Show(itm); err != nil {
		log.Printf("Couldn't find item: %d \n%+v\n", number, err)
		http.Error(w, "Couldn't find item: "+string(number), http.StatusInternalServerError)
		return
	}

	if err := jsonapi.MarshalPayload(w, itm); err != nil {
		log.Println("Couldn't unmarshal: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
