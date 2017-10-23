package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

var generatetokens = flag.Int("generatetokens", 0, `generate tokens for testing
	To generate token, developer mode should be enabled.
	To enable: export KAARYASTHAN_DEVELOPER_MODE=true`)

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

func runApp() {
	addr := config.Config.HTTPAddress
	n, _, _ := route.Router()
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
}

func generateTokens() {
	defer handleExit()
	if !config.Config.DeveloperMode {
		log.Println("Developer mode not enabled.\nTo enable: export KAARYASTHAN_DEVELOPER_MODE=true")
		panic(Exit{1})
	}

	secretKey := []byte(config.Config.TokenSecretKey)
	for i := 1; i <= *generatetokens; i++ {
		username := strconv.Itoa(i)
		var id string
		err := db.DB.QueryRow(`SELECT id FROM "users"
		WHERE username=$1 AND active=true AND email_verified=true`,
			fmt.Sprintf("developer%[1]s", username)).Scan(&id)
		if err != nil {
			if id == "" {
				*createuser = fmt.Sprintf("developer%[1]s:developer%[1]s:admin:developer%[1]s@example.org", username)
				createUser()
				err := db.DB.QueryRow(`SELECT id FROM "users"
		WHERE username=$1 AND active=true AND email_verified=true`,
					fmt.Sprintf("developer%[1]s", username)).Scan(&id)
				if err != nil {
					log.Println("Error getting user ID", err)
					panic(Exit{1})
				}
			}
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": id,
			"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		})

		tokenString, _ := token.SignedString(secretKey)
		fmt.Println(fmt.Sprintf("Token for the user 'developer%s':\n", username))
		fmt.Println(tokenString)
		fmt.Println("")
	}
}

func createUser() {
	var err error
	defer handleExit()

	args := strings.SplitN(*createuser, ":", 4)
	if len(args) != 4 {
		log.Println("Not enough segments in the string:", args)
		panic(Exit{1})
	}
	username := args[0]
	password := args[1]
	role := args[2]
	email := args[3]

	usrDS := user.NewDatastore(db.DB)
	usr := &user.User{
		Username: username,
		Name:     username,
		Email:    email,
		Role:     role,
		Password: password,
	}

	err = usrDS.Create(usr)
	if err != nil {
		log.Println("User creation failed.", err.Error())
		return
	}

	_, err = db.DB.Exec("UPDATE users SET active=true, email_verified=true, user_role=$1 WHERE id=$2",
		usr.Role, usr.ID)
	if err != nil {
		log.Println("User creation failed.", err.Error())
	}
}

func main() {
	flag.Parse()

	if *migrate {
		migrateDatabase()
		if err := db.DB.Close(); err != nil {
			log.Println("Error closing the database connection:", err)
		}
		os.Exit(0)
	}

	if *createuser != "" {
		createUser()
		if err := db.DB.Close(); err != nil {
			log.Println("Error closing the database connection:", err)
		}
		os.Exit(0)
	}

	if *generatetokens != 0 {
		generateTokens()
		if err := db.DB.Close(); err != nil {
			log.Println("Error closing the database connection:", err)
		}
		os.Exit(0)
	}

	runApp()
}
