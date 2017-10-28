package controller

import (
	"strconv"
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

func TestProjectList(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")
	defer db.DB.Exec("DELETE FROM projects")

	usrDS := user.NewDatastore(db.DB)
	usr := &user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usrDS.Create(usr); err != nil {
		t.Fatal(err)
	}

	prjDS := project.NewDatastore(db.DB)
	for i := 0; i < 5; i++ {
		prj := &project.Project{Name: "somename" + strconv.Itoa(i), Description: "Some description " + strconv.Itoa(i)}
		if err := prjDS.Create(usr, prj); err != nil {
			t.Fatal(err)
		}
		if prj.ID <= 0 {
			t.Fatalf("Data not inserted. ID: %#v", prj.ID)
		}
	}
	objs, err := prjDS.List(true)
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

		if prj.Archived {
			t.Error("Wrong Archived", prj.Archived)
		}
	}
}
