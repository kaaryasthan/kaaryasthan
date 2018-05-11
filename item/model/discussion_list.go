package item

import (
	"log"

	"github.com/pkg/errors"
)

// List searches for comments under an item
func (ds *CommentDatastore) List(itmID int) ([]*Comment, error) {
	rows, err := ds.db.Query(`SELECT id, body FROM "comments"
	WHERE deleted_at IS NULL AND item_id=$1 ORDER BY created_at`, itmID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to select all comments with item number: %d", itmID)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			log.Println(errors.Wrap(err, "failed to close database rows"))
		}
	}()

	var objs []*Comment
	for rows.Next() {
		cmt := Comment{ItemID: itmID}
		err := rows.Scan(&cmt.ID, &cmt.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		objs = append(objs, &cmt)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "failed to iterate rows")
	}
	return objs, nil
}
