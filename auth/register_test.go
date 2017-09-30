package auth

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
)

func TestUserRegister(t *testing.T) {
	defer db.DB.Exec("DELETE FROM members")
	s := Schema{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	err := s.register()
	if err != nil {
		t.Error(err)
	}
}
