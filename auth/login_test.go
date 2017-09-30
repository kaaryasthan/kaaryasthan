package auth

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
)

func TestUserLogin(t *testing.T) {
	defer db.DB.Exec("DELETE FROM members")
	s1 := Schema{Username: "jack", Password: "Secret@123"}
	err := s1.login()
	if err == nil {
		t.Error("Login succeeded")
	}
	s2 := Schema{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	err = s2.register()
	err = s1.login()
	if err == nil {
		t.Error("Login succeeded")
	}
	db.DB.Exec("UPDATE members SET active=true")
	err = s1.login()
	if err == nil {
		t.Error("Login succeeded")
	}
	db.DB.Exec("UPDATE members SET active=false, email_verified=true")
	err = s1.login()
	if err == nil {
		t.Error("Login succeeded")
	}
	db.DB.Exec("UPDATE members SET active=true, email_verified=true")
	err = s1.login()
	if err != nil {
		t.Error("Login failed", err)
	}
}
