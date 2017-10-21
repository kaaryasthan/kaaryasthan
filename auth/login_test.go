package auth

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func TestUserLogin(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")
	login := Login{Username: "jack", Password: "Secret@123"}
	if err := login.login(); err == nil {
		t.Error("Login succeeded")
	}

	usr := user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usr.Create(); err != nil {
		t.Fatal(err)
	}

	if err := login.login(); err == nil {
		t.Error("Login succeeded")
	}

	db.DB.Exec("UPDATE users SET active=true")
	if err := login.login(); err == nil {
		t.Error("Login succeeded")
	}

	db.DB.Exec("UPDATE users SET active=false, email_verified=true")
	if err := login.login(); err == nil {
		t.Error("Login succeeded")
	}

	db.DB.Exec("UPDATE users SET active=true, email_verified=true")
	if err := login.login(); err != nil {
		t.Error("Login failed", err)
	}
	if login.ID == "" {
		t.Error("Login ID not set")
	}
}
