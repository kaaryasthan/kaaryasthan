package item

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/jsonapi"
	"github.com/kaaryasthan/kaaryasthan/config"
	"github.com/kaaryasthan/kaaryasthan/search"
	"github.com/kaaryasthan/kaaryasthan/user/model"
	"github.com/lib/pq"
)

// Repository to manage items
type Repository interface {
	Create(usr *user.User, itm *Item) error
	Valid(itm *Item) error
	Show(itm *Item) error
	List(query string, offset, limit int) ([]*Item, error)
}

// Item represents an item
type Item struct {
	ID               int     `jsonapi:"primary,items"`
	Title            string  `jsonapi:"attr,title"`
	Description      string  `jsonapi:"attr,description"`
	Number           string  `jsonapi:"attr,num"`
	ProjectID        string  `jsonapi:"attr,project_id"`
	OpenState        bool    `jsonapi:"attr,open_state"`
	LockConversation bool    `jsonapi:"attr,lock_conversation"`
	CreatedBy        string  `jsonapi:"attr,created_by"`
	UpdatedBy        *string `jsonapi:"attr,updated_by"`
	Assignees        pq.StringArray
	Subscribers      pq.StringArray
	Labels           pq.StringArray
	CreatedAt        time.Time
	UpdatedAt        *time.Time
	Discussions      []*Discussion `jsonapi:"relation,discussions"`
}

// JSONAPIRelationshipLinks invoked for each relationship defined
// on the Item struct when marshaled
func (itm Item) JSONAPIRelationshipLinks(relation string) *jsonapi.Links {
	if relation == "discussions" {
		return &jsonapi.Links{
			"self":    fmt.Sprintf(config.Config.BaseURL+"/api/v1/items/%s/relationships/discussions", itm.Number),
			"related": fmt.Sprintf(config.Config.BaseURL+"/api/v1/items/%s/discussions", itm.Number),
		}
	}
	return nil
}

// Datastore implements the Repository interface
type Datastore struct {
	db *sql.DB
	bi *search.BleveIndex
}

// NewDatastore constructs a new Repository
func NewDatastore(db *sql.DB, bi *search.BleveIndex) *Datastore {
	return &Datastore{db, bi}
}

// Valid checks the validity of the item
func (ds *Datastore) Valid(itm *Item) error {
	var count int
	err := ds.db.QueryRow(`SELECT count(1) FROM "items"
		WHERE id=$1 AND lock_conversation=false AND deleted_at IS NULL`,
		itm.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Invalid item")
	}
	err = ds.db.QueryRow(`SELECT count(1) FROM "items"
		INNER JOIN projects ON items.project_id = projects.id
		WHERE items.id=$1 AND projects.archived=false AND projects.deleted_at IS NULL`,
		itm.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Invalid item")
	}
	return nil
}
