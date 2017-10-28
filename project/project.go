package controller

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// Controller holds DB
type Controller struct {
	ds  project.Repository
	uds user.Repository
}

// NewController constructs a controller
func NewController(userRepo user.Repository, repo project.Repository) *Controller {
	return &Controller{ds: repo, uds: userRepo}
}

// Register handlers
func Register(art *mux.Router, db *sql.DB) {
	c := NewController(user.NewDatastore(db), project.NewDatastore(db))

	art.HandleFunc("/api/v1/projects", c.CreateHandler).Methods("POST")
	art.HandleFunc("/api/v1/projects/{name}", c.ShowHandler).Methods("GET")
	art.HandleFunc("/api/v1/projects", c.ListHandler).Methods("GET")
}
