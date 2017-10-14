package project

import (
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func createHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	obj := new(Project)
	if err := jsonapi.UnmarshalPayload(r.Body, obj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusCreated)

	u := user.User{ID: userID}
	if err := u.Valid(); err != nil {
		log.Println("Couldn't validate user: ", err)
		return
	}

	if err := obj.Create(u); err != nil {
		log.Println("Unable to save data: ", err)
		return
	}

	if err := jsonapi.MarshalPayload(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Create creates a new project
func (obj *Project) Create(usr user.User) error {
	err := db.DB.QueryRow(`INSERT INTO "projects" (name, description, created_by) VALUES ($1, $2, $3) RETURNING id`,
		obj.Name, obj.Description, usr.ID).Scan(&obj.ID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
