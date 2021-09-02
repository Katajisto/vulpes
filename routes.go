package main

import "github.com/gorilla/mux"

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/", homeHandler)
}
