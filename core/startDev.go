//go:build !prod
// +build !prod

package core

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Startup() {
	// Create a new router
	r := mux.NewRouter()
	registerRoutes(r)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("s3/static"))))

	// Bind the router to a port
	err := http.ListenAndServe(":1337", r)
	log.Println("Exited due to error: ", err)
}
