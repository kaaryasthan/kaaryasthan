package project

import (
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/user"
)

// Show a project
func (obj *Project) Show() error {
	err := db.DB.QueryRow(`SELECT id, description, item_template, archived FROM "projects"
		WHERE name=$1 AND archived=$2 AND deleted_at IS NULL`,
		obj.Name, obj.Archived).Scan(&obj.ID, &obj.Description, &obj.ItemTemplate, &obj.Archived)
	return err
}

func showHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusOK)

	usr := user.User{ID: userID}
	if err := usr.Valid(); err != nil {
		log.Println("Couldn't validate user: "+usr.ID, err)
		http.Error(w, "Couldn't validate user: "+usr.ID, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	name := vars["name"]

	obj := &Project{Name: name}
	if err := obj.Show(); err != nil {
		log.Println("Couldn't find project: "+name, err)
		http.Error(w, "Couldn't find project: "+name, http.StatusInternalServerError)
		return
	}
	if err := jsonapi.MarshalPayload(w, obj); err != nil {
		log.Println("Couldn't unmarshal: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
