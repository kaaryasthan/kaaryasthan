package milestone

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/user"
	"github.com/lib/pq"
)

// Repository to manage labels
type Repository interface {
	Create(usr *user.User, mil *Milestone) error
	Show(mil *Milestone) error
}

// Controller holds DB
type Controller struct {
	ds  Repository
	pds project.Repository
	uds user.Repository
}

// Milestone represents a milestone
type Milestone struct {
	ID          int    `jsonapi:"primary,milestones"`
	Name        string `jsonapi:"attr,name"`
	Description string `jsonapi:"attr,description"`
	ProjectID   int    `jsonapi:"attr,project_id"`
	Items       pq.Int64Array
}

// Datastore implements the Repository interface
type Datastore struct {
	db *sql.DB
}

// NewDatastore constructs a new Repository
func NewDatastore(db *sql.DB) *Datastore {
	return &Datastore{db}
}

// NewController constructs a controller
func NewController(userRepo user.Repository, prjRepo project.Repository, repo Repository) *Controller {
	return &Controller{ds: repo, pds: prjRepo, uds: userRepo}
}

// Register handlers
func Register(art *mux.Router, db *sql.DB) {
	c := NewController(user.NewDatastore(db), project.NewDatastore(db), NewDatastore(db))

	art.HandleFunc("/api/v1/milestones", c.CreateHandler).Methods("POST")
	art.HandleFunc("/api/v1/projects/{project}/milestones/{name}", c.ShowHandler).Methods("GET")
}
