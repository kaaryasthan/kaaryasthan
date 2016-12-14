package db

import (
	"database/sql"
	"log"

	"github.com/baijum/pgmigration"
	"github.com/kaaryasthan/kaaryasthan/config"
	// DB is actually initialized here
	_ "github.com/lib/pq"
)

// DB is the database connection wrapper
var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("postgres", config.Config.PostgresConfig())
	if err != nil {
		log.Fatal(err.Error())
	}
}

// SchemaMigrate migrate database schema
func SchemaMigrate() error {
	ms := pgmigration.NewMigrationsSource(AssetNames, Asset)
	var err error
	pg, err := pgmigration.Run(DB, ms)
	if err != nil {
		return err
	}
	err = pg.Migrate("unique-code-migrations-name-00001", func(tx *sql.Tx) error { return nil })
	if err != nil {
		return err
	}
	return err
}
