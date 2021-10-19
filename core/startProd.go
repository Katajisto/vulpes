// +build prod

package core

import (
	"net/http"

	"github.com/gorilla/mux"
)

// StartProd starts the API Gateway proxy in production mode.
func Startup() {
	// Create a new router
	r := mux.NewRouter()
	registerRoutes(r)

	// Serve static assets from the static directory. In production serve these from S3 through API Gateway.
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("s3/static"))))

	// Bind the router to a port
	http.ListenAndServe(":1337", r)
}
