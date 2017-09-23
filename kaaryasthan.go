package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/kaaryasthan/kaaryasthan/auth"
	_ "github.com/kaaryasthan/kaaryasthan/auth/google"
	_ "github.com/kaaryasthan/kaaryasthan/comment"
	"github.com/kaaryasthan/kaaryasthan/config"
	"github.com/kaaryasthan/kaaryasthan/db"
	_ "github.com/kaaryasthan/kaaryasthan/item"
	"github.com/kaaryasthan/kaaryasthan/middleware"
	_ "github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/route"
	_ "github.com/kaaryasthan/kaaryasthan/web"
	"github.com/urfave/negroni"
)

//go:generate go-bindata -pkg db -o db/bindata.go -nocompress db/migrations/

var migrate = flag.Bool("migrate", false, "perform db migrations")

// Exit code for clean exit
type Exit struct {
	Code int
}

// exit code handler
func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok == true {
			os.Exit(exit.Code)
		}
		panic(e) // not an Exit, bubble up
	}
}

func init() {
	flag.Parse()
	go func() {
		defer handleExit()
		defer db.DB.Close()
		time.Sleep(5 * time.Second)
		err := db.DB.Ping()
		if err != nil {
			log.Fatal(err.Error())
		}
		if *migrate {
			err = db.SchemaMigrate()
			if err != nil {
				log.Println(err.Error())
				panic(Exit{1})
			}
			log.Println("Migration completed. Program is exiting.")
			panic(Exit{0})
		}
	}()
}

func main() {
	route.URT.PathPrefix("/api").Handler(
		negroni.New(negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext), negroni.Wrap(route.RT)))
	middleware.Run(config.Config.HTTPAddress)
}
