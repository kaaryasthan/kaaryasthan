package controller

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/label/model"
	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// Controller holds DB
type Controller struct {
	ds  label.Repository
	pds project.Repository
	uds user.Repository
}

// NewController constructs a controller
func NewController(userRepo user.Repository, prjRepo project.Repository, repo label.Repository) *Controller {
	return &Controller{ds: repo, pds: prjRepo, uds: userRepo}
}

// Register handlers
func Register(art *mux.Router, db *sql.DB) {
	c := NewController(user.NewDatastore(db), project.NewDatastore(db), label.NewDatastore(db))

	art.HandleFunc("/api/v1/labels", c.CreateHandler).Methods("POST")
}
