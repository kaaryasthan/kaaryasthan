package item

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/user"
)

func TestCommentCreate(t *testing.T) {
	defer db.DB.Exec("DELETE FROM users")
	defer db.DB.Exec("DELETE FROM projects")
	defer db.DB.Exec("DELETE FROM items")
	defer db.DB.Exec("DELETE FROM item_discussion_comment_search")
	defer db.DB.Exec("DELETE FROM discussions")
	defer db.DB.Exec("DELETE FROM comments")

	usrDS := user.NewDatastore(db.DB)
	usr := &user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usrDS.Create(usr); err != nil {
		t.Fatal(err)
	}

	prjDS := project.NewDatastore(db.DB)
	prj := &project.Project{Name: "somename", Description: "Some description"}
	if err := prjDS.Create(usr, prj); err != nil {
		t.Fatal(err)
	}

	itmDS := NewDatastore(db.DB)
	itm := &Item{Title: "sometitle", Description: "Some description", ProjectID: prj.ID}
	if err := itmDS.Create(usr, itm); err != nil {
		t.Fatal(err)
	}
	if itm.ID <= 0 {
		t.Fatalf("Data not inserted. ID: %#v", itm.ID)
	}
	if itm.Number != 1 {
		t.Fatalf("Data not inserted. Num: %#v", itm.Number)
	}

	discDS := NewDiscussionDatastore(db.DB)
	disc := &Discussion{Body: "some discussion", ItemID: itm.ID}
	if err := discDS.Create(usr, disc); err != nil {
		t.Fatal(err)
	}
	if disc.ID == "" {
		t.Fatalf("Data not inserted. ID: %#v", disc.ID)
	}

	cmtDS := NewCommentDatastore(db.DB)
	cmt := &Comment{Body: "some discussion", DiscussionID: disc.ID}
	if err := cmtDS.Create(usr, cmt); err != nil {
		t.Error(err)
	}
	if cmt.ID == "" {
		t.Errorf("Data not inserted. ID: %#v", cmt.ID)
	}
}
