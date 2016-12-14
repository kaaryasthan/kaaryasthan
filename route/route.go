package route

import "github.com/gorilla/mux"

// RT is the common router
var RT *mux.Router

func init() {
	RT = mux.NewRouter()
}
