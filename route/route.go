package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/comment"
	"github.com/kaaryasthan/kaaryasthan/item"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/web"
	"github.com/urfave/negroni"
)

// Router creates all routes
func Router() (n *negroni.Negroni, art *mux.Router, urt *mux.Router) {
	art = mux.NewRouter()
	urt = mux.NewRouter()

	urt.HandleFunc("/api/v1/projects2", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	}).Methods("GET")
	auth.Register(art, urt)
	project.Register(art, urt)
	item.Register(art, urt)
	comment.Register(art, urt)

	urt.PathPrefix("/api").Handler(
		negroni.New(negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext), negroni.Wrap(art)))
	n = negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.NewStatic(web.AssetFS()))
	n.UseHandler(urt)
	return
}
