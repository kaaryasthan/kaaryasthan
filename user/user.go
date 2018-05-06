package controller

import (
	"database/sql"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// Controller holds DB
type Controller struct {
	ds user.Repository
}

// ShowUserHandler get one user
func (c *Controller) ShowUserHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	w.Header().Set("Content-Type", jsonapi.MediaType)

	usr := &user.User{ID: userID}
	if err := c.ds.Valid(usr); err != nil {
		log.Println("Couldn't validate user: "+usr.ID, err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	username := vars["username"]

	usr = &user.User{Username: username}
	if err := c.ds.Show(usr); err != nil {
		log.Println("Couldn't find user: "+username, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err := jsonapi.MarshalPayload(w, usr); err != nil {
		log.Println("Couldn't unmarshal: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// NewController constructs a controller
func NewController(repo user.Repository) *Controller {
	return &Controller{ds: repo}
}

// Register handlers
func Register(art *mux.Router, db *sql.DB) {
	c := NewController(user.NewDatastore(db))

	// art.HandleFunc("/api/v1/users", listUsersHandler).Methods("GET")
	art.HandleFunc("/api/v1/users/{username}", c.ShowUserHandler).Methods("GET")
	//art.HandleFunc("/api/v1/users/{username}", updateUserHandler).Methods("PATCH")
	//art.HandleFunc("/api/v1/users/{username}", deleteUserHandler).Methods("DELETE")
}
