package milestone

import (
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

// Create creates a new milestone
func (ds *Datastore) Create(usr *user.User, mil *Milestone) error {
	err := ds.db.QueryRow(`INSERT INTO "milestones" (name, description, created_by, project_id) VALUES ($1, $2, $3, $4) RETURNING id`,
		mil.Name, mil.Description, usr.ID, mil.ProjectID).Scan(&mil.ID)
	return err
}
