package label

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/project/model"
	"github.com/kaaryasthan/kaaryasthan/test"
	"github.com/kaaryasthan/kaaryasthan/user/model"
)

func TestLabelCreate(t *testing.T) {
	t.Parallel()
	DB, conf := test.NewTestDB()
	defer test.ResetDB(DB, conf)

	usrDS := user.NewDatastore(DB)
	usr := &user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usrDS.Create(usr); err != nil {
		t.Fatal(err)
	}

	prjDS := project.NewDatastore(DB)
	prj := &project.Project{Name: "somename", Description: "Some description"}
	if err := prjDS.Create(usr, prj); err != nil {
		t.Error(err)
	}

	lblDS := NewDatastore(DB)
	lbl := &Label{Name: "somename", Color: "#ee0701", ProjectID: prj.ID}
	if err := lblDS.Create(usr, lbl); err != nil {
		t.Error(err)
	}
	if lbl.ID <= 0 {
		t.Errorf("Data not inserted. ID: %#v", lbl.ID)
	}
}
