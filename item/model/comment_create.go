package item

import "github.com/kaaryasthan/kaaryasthan/user/model"

// Create creates new comments
func (ds *CommentDatastore) Create(usr *user.User, cmt *Comment) error {
	err := ds.db.QueryRow(`INSERT INTO "comments" (body, created_by, discussion_id) VALUES ($1, $2, $3) RETURNING id`,
		cmt.Body, usr.ID, cmt.DiscussionID).Scan(&cmt.ID)
	return err
}
