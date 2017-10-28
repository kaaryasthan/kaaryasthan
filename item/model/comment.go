package item

import (
	"database/sql"

	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// CommentRepository to manage items
type CommentRepository interface {
	Create(usr *user.User, cmt *Comment) error
}

// Comment represents a comment
type Comment struct {
	ID           string `jsonapi:"primary,comments"`
	Body         string `jsonapi:"attr,body"`
	DiscussionID string `jsonapi:"attr,discussion_id"`
}

// CommentDatastore implements the Repository interface
type CommentDatastore struct {
	db *sql.DB
}

// NewCommentDatastore constructs a new Repository
func NewCommentDatastore(db *sql.DB) *CommentDatastore {
	return &CommentDatastore{db}
}
