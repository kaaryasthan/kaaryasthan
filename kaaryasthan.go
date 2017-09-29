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

	"github.com/jpillora/backoff"
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
	if *migrate {
		go func() {
			var err error
			defer handleExit()
			defer db.DB.Close()

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

			_, err = db.DB.Exec("SELECT 1") // db.DB.Ping() seems to be not working always
			if err != nil {
				log.Println("Migration failed.", err.Error())
				panic(Exit{1})
			}

			err = db.SchemaMigrate()
			if err != nil {
				log.Println(err.Error())
				panic(Exit{1})
			}
			log.Println("Migration completed.")
			panic(Exit{0})
		}()
	}

	n, _, _ := route.Router()
	run(config.Config.HTTPAddress, n)
}
