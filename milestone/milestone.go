package milestone

import (
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

// Milestone represents a milestone
type Milestone struct {
	ID          int    `jsonapi:"primary,milestones"`
	Name        string `jsonapi:"attr,name"`
	Description string `jsonapi:"attr,description"`
	ProjectID   int    `jsonapi:"attr,project_id"`
	Items       pq.Int64Array
}

// Register handlers
func Register(art *mux.Router) {
	art.HandleFunc("/api/v1/milestones", createHandler).Methods("POST")
	art.HandleFunc("/api/v1/projects/{project}/milestones/{name}", showHandler).Methods("GET")
}
