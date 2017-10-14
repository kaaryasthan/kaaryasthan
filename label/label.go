package label

import (
	"github.com/gorilla/mux"
)

// Label represents a label
type Label struct {
	ID        int    `jsonapi:"primary,labels"`
	Name      string `jsonapi:"attr,name"`
	Color     string `jsonapi:"attr,color"`
	ProjectID int    `jsonapi:"attr,project_id"`
}

// Register handlers
func Register(art, urt *mux.Router) {
	art.HandleFunc("/api/v1/labels", createHandler).Methods("POST")
}
