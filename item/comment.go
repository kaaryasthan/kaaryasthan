package controller

import (
	"github.com/kaaryasthan/kaaryasthan/item/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// CommentController holds DB
type CommentController struct {
	ds  item.CommentRepository
	dds item.DiscussionRepository
	uds user.Repository
}

// NewCommentController constructs a controller
func NewCommentController(userRepo user.Repository, discRepo item.DiscussionRepository, repo item.CommentRepository) *CommentController {
	return &CommentController{ds: repo, dds: discRepo, uds: userRepo}
}
