package label

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
)

func TestLabelCreate(t *testing.T) {
	defer db.DB.Exec("DELETE FROM labels")
	s := Schema{Name: "somename", Color: "#ee0701"}
	id, err := s.create()
	if err != nil {
		t.Error(err)
	}
	if id <= 0 {
		t.Errorf("Data not inserted. ID: %#v", id)
	}
}
