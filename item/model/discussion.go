package item

import (
	"database/sql"
	"errors"

	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// CommentRepository to manage items
type CommentRepository interface {
	Create(usr *user.User, cmt *Comment) error
	Valid(itm *Comment) error
	List(itm int) ([]*Comment, error)
}

// Comment represents a comment
type Comment struct {
	ID     string `jsonapi:"primary,comments"`
	Body   string `jsonapi:"attr,body"`
	ItemID int    `jsonapi:"attr,item_id"`
}

// CommentDatastore implements the Repository interface
type CommentDatastore struct {
	db *sql.DB
}

// NewCommentDatastore constructs a new Repository
func NewCommentDatastore(db *sql.DB) *CommentDatastore {
	return &CommentDatastore{db}
}

// Valid checks the validity of the comment
func (ds *CommentDatastore) Valid(cmt *Comment) error {
	var count int
	err := ds.db.QueryRow(`SELECT count(1) FROM "comments"
		WHERE id=$1 AND deleted_at IS NULL`,
		cmt.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Invalid item")
	}
	err = ds.db.QueryRow(`SELECT count(1) FROM "comments"
		INNER JOIN items ON comments.item_id = items.id
		INNER JOIN projects ON items.project_id = projects.id
		WHERE comments.id=$1 AND items.lock_conversation=false
		AND items.deleted_at IS NULL AND projects.archived=false
		AND projects.deleted_at IS NULL`,
		cmt.ID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Invalid comment")
	}
	return nil
}
