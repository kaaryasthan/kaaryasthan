package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/baijum/pgmigration"
	"github.com/jpillora/backoff"
	// DB is actually initialized here
	_ "github.com/lib/pq"
)

// Connect to database
func Connect(conf string) *sql.DB {
	var err error
	DB, err := sql.Open("postgres", conf)
	if err != nil {
		log.Fatal(err)
	}

	_, err = DB.Exec("SELECT 1") // DB.Ping() seems to be not working always

	b := &backoff.Backoff{
		Min:    7 * time.Second,
		Factor: 2,
		Max:    7 * time.Minute,
	}

	go func() {
		for {
			_, err = DB.Exec("SELECT 1") // DB.Ping() seems to be not working always
			if err != nil {
				d := b.Duration()
				log.Printf("%s (pinging failed), reconnecting in %s", err, d)
				time.Sleep(d)
				continue
			}
			b.Reset()
			time.Sleep(b.Max)
		}
	}()
	return DB
}

// SchemaMigrate migrate database schema
func SchemaMigrate(DB *sql.DB) error {
	return pgmigration.Migrate(DB, AssetNames, Asset, nil)
}
