package item

import (
	"errors"
	"time"

	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/lib/pq"
)

// Item represents an item
type Item struct {
	ID          int     `jsonapi:"primary,items"`
	Title       string  `jsonapi:"attr,title"`
	Description string  `jsonapi:"attr,description"`
	Number      int     `jsonapi:"attr,num"`
	ProjectID   int     `jsonapi:"attr,project_id"`
	OpenState   bool    `jsonapi:"attr,open_state"`
	Disabled    bool    `jsonapi:"attr,disabled"`
	CreatedBy   string  `jsonapi:"attr,created_by"`
	UpdatedBy   *string `jsonapi:"attr,updated_by"`
	Assignees   pq.StringArray
	Subscribers pq.StringArray
	Labels      pq.StringArray
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

// Valid checks the validity of the item
func (obj *Item) Valid() error {
	var count int
	err := db.DB.QueryRow(`SELECT count(1) FROM "items"
		WHERE id=$1 AND disabled=false AND deleted_at IS NULL`,
		obj.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Invalid item")
	}
	err = db.DB.QueryRow(`SELECT count(1) FROM "items"
		INNER JOIN projects ON items.project_id = projects.id
		WHERE items.id=$1 AND projects.archived=false AND projects.deleted_at IS NULL`,
		obj.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Invalid item")
	}
	return nil
}

// Register handlers
func Register(art, urt *mux.Router) {
	art.HandleFunc("/api/v1/items", createItemHandler).Methods("POST")
	art.HandleFunc("/api/v1/items/{number:[1-9]\\d*}", showItemHandler).Methods("GET").GetError()
	art.HandleFunc("/api/v1/discussions", createDiscussionHandler).Methods("POST")
	art.HandleFunc("/api/v1/comments", createCommentHandler).Methods("POST")
}
