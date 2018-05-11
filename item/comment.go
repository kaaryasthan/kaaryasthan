package controller

import (
	"github.com/kaaryasthan/kaaryasthan/item/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// CommentController holds DB
type CommentController struct {
	ds  item.CommentRepository
	ids item.Repository
	uds user.Repository
}

// NewCommentController constructs a controller
func NewCommentController(userRepo user.Repository, itmRepo item.Repository, repo item.CommentRepository) *CommentController {
	return &CommentController{ds: repo, ids: itmRepo, uds: userRepo}
}
