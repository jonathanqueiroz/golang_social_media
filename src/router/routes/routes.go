package routes

import (
	"database/sql"
	"net/http"
	"project01/src/middlewares"

	"github.com/gorilla/mux"
)

type Route struct {
	URI          string
	Method       string
	Function     func(w http.ResponseWriter, r *http.Request)
	AuthRequired bool
}

func Load(r *mux.Router, db *sql.DB) *mux.Router {
	routes := userRoutes(db)
	routes = append(routes, loginRoutes...)
	routes = append(routes, postRoutes(db)...)
	routes = append(routes, profileRoutes(db)...)
	routes = append(routes, notificationRoutes(db)...)

	for _, route := range routes {
		if route.AuthRequired {
			r.HandleFunc(route.URI, middlewares.AuthMiddleware(route.Function)).Methods(route.Method)
		} else {
			r.HandleFunc(route.URI, route.Function).Methods(route.Method)
		}
	}

	return r
}
