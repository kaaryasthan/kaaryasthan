package controller

import (
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/item/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// ListItemHandler list items
func (c *ItemController) ListItemHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)

	usr := &user.User{ID: userID}
	if err := c.uds.Valid(usr); err != nil {
		log.Println("Couldn't validate user: "+usr.ID, err)
		http.Error(w, "Couldn't validate user: "+usr.ID, http.StatusUnauthorized)
		return
	}

	vals := r.URL.Query()
	queries, ok := vals["Query"]
	var query string
	if ok {
		query = queries[0]
	}

	var objs []*item.Item
	var err error
	// FIXME: offset & limit hard-coded
	if objs, err = c.ds.List(query, 0, 20); err != nil {
		log.Println("Couldn't find items: ", err)
		http.Error(w, "Couldn't find items: ", http.StatusInternalServerError)
		return
	}

	if err := jsonapi.MarshalPayload(w, objs); err != nil {
		log.Println("Couldn't unmarshal: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
