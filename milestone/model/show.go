package milestone

// Show a milestone
func (ds *Datastore) Show(mil *Milestone) error {
	err := ds.db.QueryRow(`SELECT id, description, items FROM "milestones"
		WHERE name=$1 AND project_id=$2 AND deleted_at IS NULL`,
		mil.Name, mil.ProjectID).Scan(&mil.ID, &mil.Description, &mil.Items)
	return err
}
