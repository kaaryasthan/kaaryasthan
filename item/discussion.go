package item

import (
	"errors"

	"github.com/kaaryasthan/kaaryasthan/db"
)

// Discussion represents a discussion
type Discussion struct {
	ID     string `jsonapi:"primary,discussions"`
	Body   string `jsonapi:"attr,body"`
	ItemID int    `jsonapi:"attr,item_id"`
}

// Valid checks the validity of the discussion
func (obj *Discussion) Valid() error {
	var count int
	err := db.DB.QueryRow(`SELECT count(1) FROM "discussions"
		WHERE id=$1 AND deleted_at IS NULL`,
		obj.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Invalid item")
	}
	err = db.DB.QueryRow(`SELECT count(1) FROM "discussions"
		INNER JOIN items ON discussions.item_id = items.id
		INNER JOIN projects ON items.project_id = projects.id
		WHERE discussions.id=$1 AND items.disabled=false
		AND items.deleted_at IS NULL AND projects.archived=false
		AND projects.deleted_at IS NULL`,
		obj.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Invalid discussion")
	}
	return nil
}
