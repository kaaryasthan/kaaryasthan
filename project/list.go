package project

import (
	"database/sql"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/user"
)

// List projects
func (ds *Datastore) List(all bool) ([]Project, error) {
	var err error
	var rows *sql.Rows
	if all {
		rows, err = ds.db.Query(`SELECT id, name, description, item_template, archived FROM "projects"
		WHERE deleted_at IS NULL ORDER BY created_at`)
	} else {
		rows, err = ds.db.Query(`SELECT id, name, description, item_template, archived FROM "projects"
		WHERE archived=false AND deleted_at IS NULL ORDER BY created_at`)
	}
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Println("Error closing the database rows:", err)
		}
	}()

	var objs []Project
	for rows.Next() {
		prj := Project{}
		err = rows.Scan(&prj.ID, &prj.Name, &prj.Description, &prj.ItemTemplate, &prj.Archived)
		if err != nil {
			return nil, err
		}
		objs = append(objs, prj)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return objs, nil
}

// ListHandler list projects
func (c *Controller) ListHandler(w http.ResponseWriter, r *http.Request) {
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

	var objs []Project
	var err error
	if objs, err = c.ds.List(false); err != nil {
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
