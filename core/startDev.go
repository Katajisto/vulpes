// +build !prod

package core

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Startup() {
	// Create a new router
	r := mux.NewRouter()
	registerRoutes(r)

	// Serve static assets from the static directory. In production serve these from S3 through API Gateway.
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Bind the router to a port
	http.ListenAndServe(":8080", r)
}
