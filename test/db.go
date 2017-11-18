package test

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"hash/adler32"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kaaryasthan/kaaryasthan/config"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kelseyhightower/envconfig"
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
func NewTestDB() (*sql.DB, config.Configuration) {
	var err error
	dbname := randomDatabaseName()
	connConf := config.Config.PostgresConfig()
	DB := db.Connect(connConf)
	_, err = DB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbname))
	if err != nil {
		log.Printf("Database creation failed: %s. %#v", dbname, err)
	}
	if err = DB.Close(); err != nil {
		log.Println("Error closing the database connection:", err)
	}

	var localConfig config.Configuration
	err = envconfig.Process("kaaryasthan", &localConfig)
	if err != nil {
		log.Fatal(err.Error())
	}
	localConfig.SetDatabaseName(dbname)
	localConnConf := localConfig.PostgresConfig()
	tmpDB := db.Connect(localConnConf)
	err = db.SchemaMigrate(tmpDB)
	if err != nil {
		log.Println("Migration failed.", err.Error())
	}

	dn := strconv.Itoa(int(adler32.Checksum([]byte(dbname)))) + ".bleve"
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		log.Println(err)
	}
	bleveDir := filepath.Join(dir, dn)
	localConfig.SetBleveIndexPath(bleveDir)

	return tmpDB, localConfig
}

// ResetDB reset database to postgres
func ResetDB(DB *sql.DB, conf config.Configuration) {
	var err error
	var dbname string
	err = DB.QueryRow("SELECT current_database() as dbname").Scan(&dbname)
	if err != nil {
		log.Println("Error database name:", err)
	}

	if err = DB.Close(); err != nil {
		log.Println("Error closing the database connection:", err)
	}
	connConf := config.Config.PostgresConfig()
	baseDB := db.Connect(connConf)
	_, err = baseDB.Exec(fmt.Sprintf(`DROP DATABASE "%s"`, dbname))

	if err != nil {
		log.Printf("Database drop failed: %s. %#v", dbname, err)
	}

	err = os.RemoveAll(filepath.Dir(conf.BleveIndexPath))
	if err != nil {
		log.Printf("Remove directory failed: %s. %#v", conf.BleveIndexPath, err)
	}
}
