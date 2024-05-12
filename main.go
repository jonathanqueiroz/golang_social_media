package main

import (
	"net/http"
	"project01/src/router"
)

func main() {
	r := router.New()

	http.ListenAndServe(":8080", r)
}
