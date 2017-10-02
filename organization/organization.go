package organization

import (
	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/jsonapi"
)

// Schema represents a database schema
type Schema struct {
	Name        string
	Description string
}

// New returns a schema
func New(d jsonapi.Data) *Schema {
	s := &Schema{}
	s.Name = d.Attributes["name"]
	s.Description = d.Attributes["description"]
	return s
}

// Register handlers
func Register(art, urt *mux.Router) {
	art.HandleFunc("/api/v1/organizations", createHandler).Methods("POST")
}
