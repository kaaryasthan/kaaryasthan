package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/jpillora/backoff"
	"github.com/kaaryasthan/kaaryasthan/config"
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/route"
	"github.com/kaaryasthan/kaaryasthan/user"
)

//go:generate go-bindata -pkg db -o db/bindata.go -nocompress db/migrations/

var migrate = flag.Bool("migrate", false, "perform db migrations")
var createuser = flag.String("createuser", "", `create an active user
        Format: username:password:role:email
	e.g., admin:admin:admin:admin@example.org
	Note: username & email should not exist in the system`)

// Exit code for clean exit
type Exit struct {
	Code int
}

// exit code handler
func handleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok {
			os.Exit(exit.Code)
		}
		panic(e) // not an Exit, bubble up
	}
}

// run starts the server
func run(addr string, n http.Handler) {
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

		if err := srv.Shutdown(ctx); err != nil {
			l.Println("Error shutting down the server:", err)
		}

		if err := db.DB.Close(); err != nil {
			l.Println("Error closing DB:", err)
		}

		done <- true
	}()

	l.Printf("Listening on: %s", addr)
	<-done
}

func migrateDatabase() {
	var err error
	defer handleExit()
	defer func() {
		if err = db.DB.Close(); err != nil {
			log.Println("Error closing the database connection:", err)
		}
	}()

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
		log.Println("Migration failed.", err.Error())
		panic(Exit{1})
	}
	log.Println("Migration completed.")
	panic(Exit{0})
}

func createUser() {
	var err error
	defer handleExit()
	defer func() {
		if err = db.DB.Close(); err != nil {
			log.Println("Error closing the database connection:", err)
		}
	}()

	args := strings.SplitN(*createuser, ":", 4)
	if len(args) != 4 {
		log.Println("Not enough segments in the string:", args)
		panic(Exit{1})
	}
	username := args[0]
	password := args[1]
	role := args[2]
	email := args[3]

	usr := user.User{
		Username: username,
		Name:     username,
		Email:    email,
		Role:     role,
		Password: password,
	}

	err = usr.Create()
	if err != nil {
		log.Println("User creation failed.", err.Error())
		panic(Exit{1})
	}

	_, err = db.DB.Exec("UPDATE users SET active=true, email_verified=true, user_role=$1 WHERE id=$2",
		usr.Role, usr.ID)
	if err != nil {
		log.Println("User creation failed.", err.Error())
		panic(Exit{1})
	}
	panic(Exit{0})
}

func main() {
	flag.Parse()

	if *migrate {
		migrateDatabase()
	}

	if createuser != nil {
		createUser()
	}

	n, _, _ := route.Router()
	run(config.Config.HTTPAddress, n)
}
