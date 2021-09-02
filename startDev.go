// +build !prod

package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func startup() {
	// Create a new router
	r := mux.NewRouter()
	registerRoutes(r)
	// Bind the router to a port
	http.ListenAndServe(":8080", r)
}
