package test

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/jpillora/backoff"
	"github.com/kaaryasthan/kaaryasthan/config"
	"github.com/kaaryasthan/kaaryasthan/db"
)

func randomDatabaseName() string {
	s := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, s)
	if err != nil {
		log.Println(err)
	}
	return base64.RawStdEncoding.EncodeToString(s)
}

// NewTestDB initializes a  test database
func NewTestDB() string {
	var err error
	dbname := randomDatabaseName()
	_, err = db.DB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbname))
	if err != nil {
		log.Printf("Database creation failed: %s. %#v", dbname, err)
	}
	if err = db.DB.Close(); err != nil {
		log.Println("Error closing the database connection:", err)
	}

	config.Config.SetDatabaseName(dbname)
	db.DB, _ = sql.Open("postgres", config.Config.PostgresConfig())

	b := &backoff.Backoff{
		Min:    7 * time.Second,
		Factor: 2,
		Max:    7 * time.Minute,
	}

	for i := 0; i < 7; i++ {
		_, err = db.DB.Exec("SELECT 1") // db.DB.Ping() seems to be not working always
		if err != nil {
			d := b.Duration()
			log.Printf("%s (pinging failed), reconnecting in %s", err, d)
			time.Sleep(d)
			continue
		}
		b.Reset()
	}
	err = db.SchemaMigrate()
	if err != nil {
		log.Println("Migration failed.", err.Error())
	}

	return dbname
}

// ResetDB reset database to postgres
func ResetDB(dbname string) {
	var err error
	if err = db.DB.Close(); err != nil {
		log.Println("Error closing the database connection:", err)
	}

	config.Config.SetDatabaseName("postgres")
	db.DB, _ = sql.Open("postgres", config.Config.PostgresConfig())

	b := &backoff.Backoff{
		Min:    7 * time.Second,
		Factor: 2,
		Max:    7 * time.Minute,
	}

	for i := 0; i < 7; i++ {
		_, err = db.DB.Exec("SELECT 1") // db.DB.Ping() seems to be not working always
		if err != nil {
			d := b.Duration()
			log.Printf("%s (pinging failed), reconnecting in %s", err, d)
			time.Sleep(d)
			continue
		}
		b.Reset()
	}

	_, err = db.DB.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, dbname))

	if err != nil {
		log.Printf("Database drop failed: %s. %#v", dbname, err)
	}
}
