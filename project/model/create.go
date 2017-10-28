package project

import "github.com/kaaryasthan/kaaryasthan/user/model"

// Create creates a new project
func (ds *Datastore) Create(usr *user.User, prj *Project) error {
	err := ds.db.QueryRow(`INSERT INTO "projects" (name, description, created_by) VALUES ($1, $2, $3) RETURNING id`,
		prj.Name, prj.Description, usr.ID).Scan(&prj.ID)
	return err
}
