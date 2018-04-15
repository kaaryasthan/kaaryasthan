package project

import (
	"database/sql"
	"log"
)

// List projects
func (ds *Datastore) List(all bool) ([]*Project, error) {
	var err error
	var rows *sql.Rows
	if all {
		rows, err = ds.db.Query(`SELECT id, name, description, item_template, archived FROM "projects"
		WHERE deleted_at IS NULL ORDER BY created_at`)
	} else {
		rows, err = ds.db.Query(`SELECT id, name, description, item_template, archived FROM "projects"
		WHERE archived=false AND deleted_at IS NULL ORDER BY created_at`)
	}
	if err != nil {
		return nil, err
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Println("Error closing the database rows:", err)
		}
	}()

	var objs []*Project
	for rows.Next() {
		prj := Project{}
		err = rows.Scan(&prj.ID, &prj.Name, &prj.Description, &prj.ItemTemplate, &prj.Archived)
		if err != nil {
			return nil, err
		}
		objs = append(objs, &prj)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return objs, nil
}
