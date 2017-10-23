package project

import (
	"database/sql"
	"errors"

	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/user"
)

// Repository to manage projects
type Repository interface {
	Create(usr *user.User, prj *Project) error
	Valid(prj *Project) error
	Show(prj *Project) error
	List(all bool) ([]Project, error)
}

// Controller holds DB
type Controller struct {
	ds  Repository
	uds user.Repository
}

// Project represents a project
type Project struct {
	ID           int    `jsonapi:"primary,projects"`
	Name         string `jsonapi:"attr,name"`
	Description  string `jsonapi:"attr,description"`
	ItemTemplate string `jsonapi:"attr,item_template"`
	Archived     bool   `jsonapi:"attr,archived"`
}

// Datastore implements the Repository interface
type Datastore struct {
	db *sql.DB
}

// NewDatastore constructs a new Repository
func NewDatastore(db *sql.DB) *Datastore {
	return &Datastore{db}
}

// Valid checks the validity of the project
func (ds *Datastore) Valid(prj *Project) error {
	var count int
	err := ds.db.QueryRow(`SELECT count(1) FROM "projects"
		WHERE id=$1 AND archived=false AND deleted_at IS NULL`,
		prj.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Invalid project")
	}
	return nil
}

// NewController constructs a controller
func NewController(userRepo user.Repository, repo Repository) *Controller {
	return &Controller{ds: repo, uds: userRepo}
}

// Register handlers
func Register(art *mux.Router, db *sql.DB) {
	c := NewController(user.NewDatastore(db), NewDatastore(db))

	art.HandleFunc("/api/v1/projects", c.CreateHandler).Methods("POST")
	art.HandleFunc("/api/v1/projects/{name}", c.ShowHandler).Methods("GET")
	art.HandleFunc("/api/v1/projects", c.ListHandler).Methods("GET")
}
