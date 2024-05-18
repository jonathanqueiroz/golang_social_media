package main

import (
	"log"
	"net/http"
	"project01/src/config"
	"project01/src/db"
	"project01/src/router"
)

func main() {
	config.Load()

	db, err := db.New()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := router.New(db)

	http.ListenAndServe(":8080", r)
}
