package router

import (
	"project01/src/router/routes"

	"github.com/gorilla/mux"
)

func New() *mux.Router {
	r := mux.NewRouter()
	routes.Load(r)
	return r
}
