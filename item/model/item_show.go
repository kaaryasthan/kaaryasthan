package item

import "github.com/pkg/errors"

// Show an item
func (ds *Datastore) Show(itm *Item) error {
	err := ds.db.QueryRow(`SELECT id, title, description, open_state,
			project_id, lock_conversation, created_by, updated_by, assignees,
			subscribers, labels, created_at, updated_at FROM "items"
			WHERE num=$1 AND deleted_at IS NULL`,
		itm.Number).Scan(&itm.ID, &itm.Title, &itm.Description,
		&itm.OpenState, &itm.ProjectID, &itm.LockConversation, &itm.CreatedBy,
		&itm.UpdatedBy, &itm.Assignees, &itm.Subscribers, &itm.Labels,
		&itm.CreatedAt, &itm.UpdatedAt)
	return errors.Wrap(err, "select item with number: "+itm.Number)
}
