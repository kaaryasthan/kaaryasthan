package controller

import (
	"github.com/kaaryasthan/kaaryasthan/item/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// DiscussionController holds DB
type DiscussionController struct {
	ds  item.DiscussionRepository
	ids item.Repository
	uds user.Repository
}

// NewDiscussionController constructs a controller
func NewDiscussionController(userRepo user.Repository, itmRepo item.Repository, repo item.DiscussionRepository) *DiscussionController {
	return &DiscussionController{ds: repo, ids: itmRepo, uds: userRepo}
}
