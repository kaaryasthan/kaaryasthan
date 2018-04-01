package item

import (
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// Create an item in the database
func (ds *Datastore) Create(usr *user.User, itm *Item) error {
	err := ds.db.QueryRow(`INSERT INTO "items" (title, description, created_by, project_id, labels) VALUES
		($1, $2, $3, $4, $5) RETURNING id, num`,
		itm.Title, itm.Description, usr.ID, itm.ProjectID, itm.Labels).Scan(&itm.ID, &itm.Number)
	return err
}
