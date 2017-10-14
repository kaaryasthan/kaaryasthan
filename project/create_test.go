package project

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func TestProjectCreate(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")
	defer db.DB.Exec("DELETE FROM projects")

	usr := user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usr.Create(); err != nil {
		t.Error(err)
	}

	prj := Project{Name: "somename", Description: "Some description"}
	if err := prj.Create(usr); err != nil {
		t.Error(err)
	}
	if prj.ID <= 0 {
		t.Errorf("Data not inserted. ID: %#v", prj.ID)
	}
}
