package organization

import (
	"testing"

	"github.com/kaaryasthan/kaaryasthan/db"
)

func TestOrganizationCreate(t *testing.T) {
	defer db.DB.Exec("DELETE FROM organizations")
	s := Schema{Name: "somename", Description: "Some description"}
	id, err := s.create()
	if err != nil {
		t.Error(err)
	}
	if id <= 0 {
		t.Errorf("Data not inserted. ID: %#v", id)
	}
}
