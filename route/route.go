package route

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaaryasthan/kaaryasthan/auth"
	"github.com/kaaryasthan/kaaryasthan/item"
	"github.com/kaaryasthan/kaaryasthan/label"
	"github.com/kaaryasthan/kaaryasthan/milestone"
	"github.com/kaaryasthan/kaaryasthan/project"
	"github.com/kaaryasthan/kaaryasthan/user"
	"github.com/kaaryasthan/kaaryasthan/web"
	"github.com/thoas/stats"
	"github.com/urfave/negroni"
)

// Router creates all routes
func Router() (n *negroni.Negroni, art *mux.Router, urt *mux.Router) {
	art = mux.NewRouter()
	urt = mux.NewRouter()

	middleware := stats.New()

	art.HandleFunc("/api/v1/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		stats := middleware.Data()

		b, _ := json.Marshal(stats)

		if _, err := w.Write(b); err != nil {
			log.Println("Couldn't write to response:", err)
		}
	})

	user.Register(art)
	auth.Register(urt)
	project.Register(art)
	milestone.Register(art)
	label.Register(art)
	item.Register(art)

	urt.PathPrefix("/api").Handler(
		negroni.New(negroni.HandlerFunc(auth.JwtMiddleware.HandlerWithNext), negroni.Wrap(art)))
	n = negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.NewStatic(web.AssetFS()))
	n.Use(middleware)
	n.UseHandler(urt)
	return
}
