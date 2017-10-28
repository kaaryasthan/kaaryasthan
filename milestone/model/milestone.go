package milestone

import (
	"database/sql"

	"github.com/kaaryasthan/kaaryasthan/user/model"
	"github.com/lib/pq"
)

// Repository to manage labels
type Repository interface {
	Create(usr *user.User, mil *Milestone) error
	Show(mil *Milestone) error
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
