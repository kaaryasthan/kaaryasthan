package auth

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func TestUserLogin(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")
	s1 := Login{Username: "jack", Password: "Secret@123"}
	err := s1.login()
	if err == nil {
		t.Error("Login succeeded")
	}
	s2 := user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	err = s2.Create()
	err = s1.login()
	if err == nil {
		t.Error("Login succeeded")
	}
	db.DB.Exec("UPDATE users SET active=true")
	err = s1.login()
	if err == nil {
		t.Error("Login succeeded")
	}
	db.DB.Exec("UPDATE users SET active=false, email_verified=true")
	err = s1.login()
	if err == nil {
		t.Error("Login succeeded")
	}
	db.DB.Exec("UPDATE users SET active=true, email_verified=true")
	err = s1.login()
	if err != nil {
		t.Error("Login failed", err)
	}
	if s1.ID == "" {
		t.Error("Login ID not set")
	}
}
