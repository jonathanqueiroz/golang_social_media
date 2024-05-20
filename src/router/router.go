package router

import (
	"database/sql"
	"project01/src/router/routes"
	"project01/src/websocket"

	"github.com/gorilla/mux"
)

func New(db *sql.DB) *mux.Router {
	r := mux.NewRouter()
	routes.Load(r, db)

	r.HandleFunc("/ws", websocket.HandleConnections)
	return r
}
