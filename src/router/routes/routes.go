package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	URI          string
	Method       string
	Function     func(w http.ResponseWriter, r *http.Request)
	AuthRequired bool
}

func Load(r *mux.Router) *mux.Router {
	routes := userRoutes
	routes = append(routes, postRoutes...)

	for _, route := range routes {
		r.HandleFunc(route.URI, route.Function).Methods(route.Method)
	}

	return r
}
