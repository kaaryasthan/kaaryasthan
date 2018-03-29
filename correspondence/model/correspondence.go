package correspondence

import "database/sql"

// Repository provides sending mail
type Repository interface {
	SendMail(eml *Email) error
}

// Datastore implements the Repository interface
type Datastore struct {
	db *sql.DB
}

// NewDatastore constructs a new Repository
func NewDatastore(db *sql.DB) *Datastore {
	return &Datastore{db}
}
