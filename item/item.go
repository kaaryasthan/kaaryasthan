package item

import (
	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/route"
)

// Data represents a item payload
type Data struct {
	Type       string            `json:"type"`
	ID         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
}

// Schema represents a database schema
type Schema struct {
	Title       string
	Description string
}

// Create creates a new item
func (obj *Schema) Create() (int, error) {
	var id int
	err := db.DB.QueryRow(`INSERT INTO "items" (title, description) VALUES ($1, $2) RETURNING id`, obj.Title, obj.Description).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// New returns a schema
func New(d Data) *Schema {
	s := &Schema{}
	s.Title = d.Attributes["title"]
	s.Description = d.Attributes["description"]
	return s
}

func init() {
	route.RT.HandleFunc("/api/v1/items", createHandler).Methods("POST")
}
