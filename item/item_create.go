package item

import (
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func createItemHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	obj := new(Item)
	if err := jsonapi.UnmarshalPayload(r.Body, obj); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusCreated)

	usr := user.User{ID: userID}
	if err := usr.Valid(); err != nil {
		log.Println("Couldn't validate user: ", err)
		return
	}

	prj := project.Project{ID: obj.ProjectID}
	if err := prj.Valid(); err != nil {
		log.Println("Couldn't validate project: ", err)
		return
	}

	err := obj.Create(usr)
	if err != nil {
		log.Println("Unable to save data: ", err)
		return
	}

	if err := jsonapi.MarshalPayload(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Create an item in the database
func (obj *Item) Create(usr user.User) error {
	err := db.DB.QueryRow(`INSERT INTO "items" (title, description, created_by, project_id) VALUES
		($1, $2, $3, $4) RETURNING id, num`,
		obj.Title, obj.Description, usr.ID, obj.ProjectID).Scan(&obj.ID, &obj.Number)
	return err
}
