package route

import "github.com/gorilla/mux"

var (
	// RT is the common router (require authentication)
	RT *mux.Router
	// URT is the router which doesn't require authentication
	URT *mux.Router
)

func init() {
	RT = mux.NewRouter()
	URT = mux.NewRouter()
}
