package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/baijum/pgmigration"
	"github.com/jpillora/backoff"
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
}

// SchemaMigrate migrate database schema
func SchemaMigrate() error {
	return pgmigration.Migrate(DB, AssetNames, Asset, nil)
}
