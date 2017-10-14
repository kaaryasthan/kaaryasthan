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

	usr := user.User{Username: "jack", Name: "Jack Wilber", Email: "jack@example.com", Password: "Secret@123"}
	if err := usr.Create(); err != nil {
		t.Error(err)
	}

	prj := project.Project{Name: "somename", Description: "Some description"}
	if err := prj.Create(usr); err != nil {
		t.Error(err)
	}

	itm := Item{Title: "sometitle", Description: "Some description", ProjectID: prj.ID}
	if err := itm.Create(usr); err != nil {
		t.Fatal(err)
	}
	if itm.ID <= 0 {
		t.Fatalf("Data not inserted. ID: %#v", itm.ID)
	}
	if itm.Num != 1 {
		t.Fatalf("Data not inserted. Num: %#v", itm.Num)
	}

	disc := Discussion{Body: "some discussion", ItemID: itm.ID}
	if err := disc.Create(usr); err != nil {
		t.Fatal(err)
	}
	if disc.ID == "" {
		t.Fatalf("Data not inserted. ID: %#v", disc.ID)
	}

	com := Comment{Body: "some discussion", DiscussionID: disc.ID}
	if err := com.Create(usr); err != nil {
		t.Error(err)
	}
	if com.ID == "" {
		t.Errorf("Data not inserted. ID: %#v", com.ID)
	}
}
