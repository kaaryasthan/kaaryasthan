package test

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"

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
func NewTestDB() *sql.DB {
	var err error
	dbname := randomDatabaseName()
	conf := config.Config.PostgresConfig()
	DB := db.Connect(conf)
	_, err = DB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbname))
	if err != nil {
		log.Printf("Database creation failed: %s. %#v", dbname, err)
	}
	if err = DB.Close(); err != nil {
		log.Println("Error closing the database connection:", err)
	}

	config.Config.SetDatabaseName(dbname)
	conf = config.Config.PostgresConfig()
	tmpDB := db.Connect(conf)
	err = db.SchemaMigrate(tmpDB)
	if err != nil {
		log.Println("Migration failed.", err.Error())
	}
	return tmpDB
}

// ResetDB reset database to postgres
func ResetDB(DB *sql.DB) {
	var err error
	var dbname string
	err = DB.QueryRow("SELECT current_database() as dbname").Scan(&dbname)
	if err != nil {
		log.Println("Error database name:", err)
	}

	if err = DB.Close(); err != nil {
		log.Println("Error closing the database connection:", err)
	}

	config.Config.SetDatabaseName("postgres")
	conf := config.Config.PostgresConfig()
	tmpDB := db.Connect(conf)
	_, err = tmpDB.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, dbname))

	if err != nil {
		log.Printf("Database drop failed: %s. %#v", dbname, err)
	}
}
