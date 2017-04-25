package comment

import (
	"log"
	"strconv"

	"github.com/kaaryasthan/kaaryasthan/db"
	"github.com/kaaryasthan/kaaryasthan/route"
)

// ItemData represents a item data
type ItemData struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Data represents a comment payload
type Data struct {
	Type          string            `json:"type"`
	ID            string            `json:"id"`
	Attributes    map[string]string `json:"attributes"`
	Relationships map[string]struct {
		Data ItemData `json:"data"`
	} `json:"relationships"`
}

// Schema represents a database schema
type Schema struct {
	Body string
	Item int
}

// Create create new comment
func (obj *Schema) Create() (int, error) {
	var id int
	err := db.DB.QueryRow(`INSERT INTO "comments" (body, item) VALUES ($1, $2) RETURNING id`, obj.Body, obj.Item).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// New returns a schema
func New(d Data) (*Schema, error) {
	s := &Schema{}
	s.Body = d.Attributes["body"]
	id, err := strconv.Atoi(d.Relationships["items"].Data.ID)
	if err != nil {
		log.Println("Unable to convert to int: ", err)
		return nil, err
	}
	s.Item = id
	return s, nil
}

func init() {
	route.RT.HandleFunc("/api/v1/items/{item}/comments", createHandler).Methods("POST")
}