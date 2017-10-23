package item

import (
	"database/sql"

	"github.com/kaaryasthan/kaaryasthan/user"
)

// CommentRepository to manage items
type CommentRepository interface {
	Create(usr *user.User, cmt *Comment) error
}

// CommentController holds DB
type CommentController struct {
	ds  CommentRepository
	dds DiscussionRepository
	uds user.Repository
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

// NewCommentController constructs a controller
func NewCommentController(userRepo user.Repository, discRepo DiscussionRepository, repo CommentRepository) *CommentController {
	return &CommentController{ds: repo, dds: discRepo, uds: userRepo}
}
