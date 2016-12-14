package project

import (
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/route"
)

// Data represents a project payload
type Data struct {
	Type       string            `json:"type"`
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
}

// Schema represents a database schema
type Schema struct {
	Name        string
	Description string
}

func (obj *Schema) create() (int, error) {
	var id int
	err := db.DB.QueryRow(`INSERT INTO "projects" (name, description) VALUES ($1, $2) RETURNING id`, obj.Name, obj.Description).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// New returns a schema
func New(d Data) *Schema {
	s := &Schema{}
	s.Name = d.Attributes["name"]
	s.Description = d.Attributes["description"]
	return s
}

func init() {
	route.RT.HandleFunc("/api/v1/projects", createHandler).Methods("POST")
}
