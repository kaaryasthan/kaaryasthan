package milestone

import (
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/user"
)

// CreateHandler creates milestone
func (c *Controller) CreateHandler(w http.ResponseWriter, r *http.Request) {
	tkn := r.Context().Value("user").(*jwt.Token)
	userID := tkn.Claims.(jwt.MapClaims)["sub"].(string)

	w.Header().Set("Content-Type", jsonapi.MediaType)
	w.WriteHeader(http.StatusCreated)

	usr := &user.User{ID: userID}
	if err := c.uds.Valid(usr); err != nil {
		log.Println("Couldn't validate user: ", err)
		return
	}

	mil := new(Milestone)
	if err := jsonapi.UnmarshalPayload(r.Body, mil); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prj := &project.Project{ID: mil.ProjectID}
	if err := c.pds.Valid(prj); err != nil {
		log.Println("Couldn't validate project: ", err)
		return
	}

	err := c.ds.Create(usr, mil)
	if err != nil {
		log.Println("Unable to save data: ", err)
		return
	}

	if err := jsonapi.MarshalPayload(w, mil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Create creates a new milestone
func (ds *Datastore) Create(usr *user.User, mil *Milestone) error {
	err := db.DB.QueryRow(`INSERT INTO "milestones" (name, description, created_by, project_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		mil.Name, mil.Description, usr.ID, mil.ProjectID).Scan(&mil.ID)
	return err
}
