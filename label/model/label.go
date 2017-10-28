package label

import (
	"database/sql"

	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// Repository to manage labels
type Repository interface {
	Create(usr *user.User, lbl *Label) error
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
