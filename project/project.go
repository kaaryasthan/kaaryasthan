package project

import (
	"errors"

	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/db"
)

// Project represents a project
type Project struct {
	ID          int    `jsonapi:"primary,projects"`
	Name        string `jsonapi:"attr,name"`
	Description string `jsonapi:"attr,description"`
}

// Valid checks the validity of the project
func (obj *Project) Valid() error {
	var count int
	err := db.DB.QueryRow(`SELECT count(1) FROM "projects"
		WHERE id=$1 AND archived=false AND deleted_at IS NULL`,
		obj.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Invalid project")
	}
	return nil
}

// Register handlers
func Register(art, urt *mux.Router) {
	art.HandleFunc("/api/v1/projects", createHandler).Methods("POST")
}
