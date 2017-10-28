package item

import "github.com/kaaryasthan/kaaryasthan/user/model"

// Create creates new discussions
func (ds *DiscussionDatastore) Create(usr *user.User, disc *Discussion) error {
	err := ds.db.QueryRow(`INSERT INTO "discussions" (body, created_by, item_id) VALUES ($1, $2, $3) RETURNING id`,
		disc.Body, usr.ID, disc.ItemID).Scan(&disc.ID)
	return err
}
