package controller

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/item/model"
	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/search"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// ItemController holds DB
type ItemController struct {
	ds  item.Repository
	pds project.Repository
	uds user.Repository
}

// NewItemController constructs a controller
func NewItemController(userRepo user.Repository, prjRepo project.Repository, repo item.Repository) *ItemController {
	return &ItemController{ds: repo, pds: prjRepo, uds: userRepo}
}

// Register handlers
func Register(art *mux.Router, db *sql.DB, bi *search.BleveIndex) {
	ic := NewItemController(user.NewDatastore(db), project.NewDatastore(db), item.NewDatastore(db, bi))
	dc := NewDiscussionController(user.NewDatastore(db), item.NewDatastore(db, bi), item.NewDiscussionDatastore(db))
	cc := NewCommentController(user.NewDatastore(db), item.NewDiscussionDatastore(db), item.NewCommentDatastore(db))

	art.HandleFunc("/api/v1/items", ic.CreateItemHandler).Methods("POST")
	art.HandleFunc("/api/v1/items/{number:[1-9]\\d*}", ic.ShowItemHandler).Methods("GET")
	art.HandleFunc("/api/v1/items", ic.ListItemHandler).Methods("GET")
	art.HandleFunc("/api/v1/discussions", dc.CreateDiscussionHandler).Methods("POST")
	art.HandleFunc("/api/v1/items/{number:[1-9]\\d*}/relationships/discussions", dc.CreateDiscussionHandler).Methods("POST")
	art.HandleFunc("/api/v1/items/{number:[1-9]\\d*}/discussions", dc.ListDiscussionHandler).Methods("GET")
	art.HandleFunc("/api/v1/discussions/{id}/relationships/comments", cc.CreateCommentHandler).Methods("POST")
}
