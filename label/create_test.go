package label

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func TestLabelCreate(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")
	defer db.DB.Exec("DELETE FROM projects")
	defer db.DB.Exec("DELETE FROM labels")

	usrDS := user.NewDatastore(db.DB)
	usr := &user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usrDS.Create(usr); err != nil {
		t.Fatal(err)
	}

	prjDS := project.NewDatastore(db.DB)
	prj := &project.Project{Name: "somename", Description: "Some description"}
	if err := prjDS.Create(usr, prj); err != nil {
		t.Error(err)
	}

	lblDS := NewDatastore(db.DB)
	lbl := &Label{Name: "somename", Color: "#ee0701", ProjectID: prj.ID}
	if err := lblDS.Create(usr, lbl); err != nil {
		t.Error(err)
	}
	if lbl.ID <= 0 {
		t.Errorf("Data not inserted. ID: %#v", lbl.ID)
	}
}
