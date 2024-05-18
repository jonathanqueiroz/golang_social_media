package router

import (
	"database/sql"
	"project01/src/router/routes"

	"github.com/gorilla/mux"
)

func New(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	routes.Load(r, db)
	return r
}
