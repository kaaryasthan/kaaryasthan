package label

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/user"
)

// Repository to manage labels
type Repository interface {
	Create(usr *user.User, lbl *Label) error
}

// Controller holds DB
type Controller struct {
	ds  Repository
	pds project.Repository
	uds user.Repository
}

// Label represents a label
type Label struct {
	ID        int    `jsonapi:"primary,labels"`
	Name      string `jsonapi:"attr,name"`
	Color     string `jsonapi:"attr,color"`
	ProjectID int    `jsonapi:"attr,project_id"`
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

	art.HandleFunc("/api/v1/labels", c.CreateHandler).Methods("POST")
}
