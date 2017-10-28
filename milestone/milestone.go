package controller

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/milestone/model"
	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// Controller holds DB
type Controller struct {
	ds  milestone.Repository
	pds project.Repository
	uds user.Repository
}

// NewController constructs a controller
func NewController(userRepo user.Repository, prjRepo project.Repository, repo milestone.Repository) *Controller {
	return &Controller{ds: repo, pds: prjRepo, uds: userRepo}
}

// Register handlers
func Register(art *mux.Router, db *sql.DB) {
	c := NewController(user.NewDatastore(db), project.NewDatastore(db), milestone.NewDatastore(db))

	art.HandleFunc("/api/v1/milestones", c.CreateHandler).Methods("POST")
	art.HandleFunc("/api/v1/projects/{project}/milestones/{name}", c.ShowHandler).Methods("GET")
}
