package main

import (
	"net/http"
	"project01/src/config"
	"project01/src/router"
)

func main() {
	config.Load()
	r := router.New()

	http.ListenAndServe(":8080", r)
}
