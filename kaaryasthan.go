package main

import (
	"flag"
	"log"
	"os"
	"time"

	_ "github.com/kaaryasthan/kaaryasthan/auth"
	_ "github.com/kaaryasthan/kaaryasthan/comment"
	"github.com/kaaryasthan/kaaryasthan/config"
	"github.com/kaaryasthan/kaaryasthan/db"
	_ "github.com/kaaryasthan/kaaryasthan/item"
	"github.com/kaaryasthan/kaaryasthan/middleware"
	_ "github.com/kaaryasthan/kaaryasthan/project"
	_ "github.com/kaaryasthan/kaaryasthan/route"
	_ "github.com/kaaryasthan/kaaryasthan/web"
)

//go:generate go-bindata -pkg db -o db/bindata.go -nocompress db/migrations/

var migrate = flag.Bool("migrate", false, "perform db migrations")

func init() {
	flag.Parse()
	go func() {
		time.Sleep(5 * time.Second)
		err := db.DB.Ping()
		if err != nil {
			log.Fatal(err.Error())
		}
		if *migrate {
			err = db.SchemaMigrate()
			if err != nil {
				log.Fatal(err.Error())
			}
			log.Println("Migration completed. Program is exiting.")
			os.Exit(0)
		}
	}()
}

func main() {
	middleware.Run(config.Config.KaaryasthanAddress)
}
