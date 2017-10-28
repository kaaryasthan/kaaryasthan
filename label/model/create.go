package label

import "github.com/kaaryasthan/kaaryasthan/user/model"

// Create creates a new label
func (ds *Datastore) Create(usr *user.User, lbl *Label) error {
	err := ds.db.QueryRow(`INSERT INTO "labels" (name, color, created_by, project_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		lbl.Name, lbl.Color, usr.ID, lbl.ProjectID).Scan(&lbl.ID)
	return err
}
