package project

import (
	"strconv"
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func TestProjectList(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")
	defer db.DB.Exec("DELETE FROM projects")

	usr := user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usr.Create(); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		prj := Project{Name: "somename" + strconv.Itoa(i), Description: "Some description " + strconv.Itoa(i)}
		if err := prj.Create(usr); err != nil {
			t.Fatal(err)
		}
		if prj.ID <= 0 {
			t.Fatalf("Data not inserted. ID: %#v", prj.ID)
		}
	}
	objs, err := List(true)
	if err != nil {
		t.Error("There are 5 projects", err)
	}

	if len(objs) != 5 {
		t.Error("There are 5 projects")
	}

	for i := 0; i < 5; i++ {
		prj := objs[i]
		if prj.Name != "somename"+strconv.Itoa(i) {
			t.Error("Wrong Name", prj.Name)
		}

		if prj.Description != "Some description "+strconv.Itoa(i) {
			t.Error("Wrong Description", prj.Description)
		}

		if prj.Archived == true {
			t.Error("Wrong Archived", prj.Archived)
		}
	}
}
