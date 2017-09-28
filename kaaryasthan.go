package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kaaryasthan/kaaryasthan/config"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/route"
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

// run starts the server
func run(addr string, n *negroni.Negroni) {
	l := log.New(os.Stdout, "[kaaryasthan] ", 0)

	stopChan := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	srv := &http.Server{Addr: addr, Handler: n}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			l.Println("Error starting server:", err)

			if err := db.DB.Close(); err != nil {
				l.Println("Error closing DB:", err)
			}

			done <- true
		}
	}()

	go func() {

		sig := <-stopChan
		log.Println("Signal received:", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)

		// Even though ctx will be expired, it is good practice to call its
		// cancelation function in any case. Failure to do so may keep the
		// context and its parent alive longer than necessary.
		defer cancel()

		srv.Shutdown(ctx)

		if err := db.DB.Close(); err != nil {
			l.Println("Error closing DB:", err)
		}

		done <- true
	}()

	l.Printf("Listening on: %s", addr)
	<-done
}

func main() {
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

	n, _, _ := route.Router()
	run(config.Config.HTTPAddress, n)
}
