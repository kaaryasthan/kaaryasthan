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

// CreateItemHandler creates item
func (c *Controller) CreateItemHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusCreated)

	usr := &user.User{ID: userID}
	if err := c.uds.Valid(usr); err != nil {
		log.Println("Couldn't validate user: ", err)
		return
	}

	itm := new(Item)
	if err := jsonapi.UnmarshalPayload(r.Body, itm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prj := &project.Project{ID: itm.ProjectID}
	if err := c.pds.Valid(prj); err != nil {
		log.Println("Couldn't validate project: ", err)
		return
	}

	err := c.ds.Create(usr, itm)
	if err != nil {
		log.Println("Unable to save data: ", err)
		return
	}

	if err := jsonapi.MarshalPayload(w, itm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Create an item in the database
func (ds *Datastore) Create(usr *user.User, itm *Item) error {
	err := db.DB.QueryRow(`INSERT INTO "items" (title, description, created_by, project_id) VALUES
		($1, $2, $3, $4) RETURNING id, num`,
		itm.Title, itm.Description, usr.ID, itm.ProjectID).Scan(&itm.ID, &itm.Number)
	return err
}
