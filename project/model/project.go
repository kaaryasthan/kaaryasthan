package project

import (
	"database/sql"
	"errors"

	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// Repository to manage projects
type Repository interface {
	Create(usr *user.User, prj *Project) error
	Valid(prj *Project) error
	Show(prj *Project) error
	List(all bool) ([]*Project, error)
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
