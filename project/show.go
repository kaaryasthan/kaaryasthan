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
func (ds *Datastore) Show(prj *Project) error {
	err := db.DB.QueryRow(`SELECT id, description, item_template, archived FROM "projects"
		WHERE name=$1 AND archived=$2 AND deleted_at IS NULL`,
		prj.Name, prj.Archived).Scan(&prj.ID, &prj.Description, &prj.ItemTemplate, &prj.Archived)
	return err
}

// ShowHandler shows project
func (c *Controller) ShowHandler(w http.ResponseWriter, r *http.Request) {
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

	vars := mux.Vars(r)
	name := vars["name"]

	prj := &Project{Name: name}
	if err := c.ds.Show(prj); err != nil {
		log.Println("Couldn't find project: "+name, err)
		http.Error(w, "Couldn't find project: "+name, http.StatusInternalServerError)
		return
	}
	if err := jsonapi.MarshalPayload(w, prj); err != nil {
		log.Println("Couldn't unmarshal: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
