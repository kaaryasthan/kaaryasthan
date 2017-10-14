package milestone

import (
	"github.com/gorilla/mux"
)

// Milestone represents a milestone
type Milestone struct {
	ID          int    `jsonapi:"primary,milestones"`
	Name        string `jsonapi:"attr,name"`
	Description string `jsonapi:"attr,description"`
	ProjectID   int    `jsonapi:"attr,project_id"`
}

// Register handlers
func Register(art, urt *mux.Router) {
	art.HandleFunc("/api/v1/milestones", createHandler).Methods("POST")
}
