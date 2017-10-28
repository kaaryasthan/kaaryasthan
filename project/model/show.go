package project

// Show a project
func (ds *Datastore) Show(prj *Project) error {
	err := ds.db.QueryRow(`SELECT id, description, item_template, archived FROM "projects"
		WHERE name=$1 AND archived=$2 AND deleted_at IS NULL`,
		prj.Name, prj.Archived).Scan(&prj.ID, &prj.Description, &prj.ItemTemplate, &prj.Archived)
	return err
}
