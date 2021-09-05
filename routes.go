package main

import (
	"fmt"
	"html"
	"net/http"

	"github.com/gorilla/mux"
)

func registerRoutes(r *mux.Router) {
	r.HandleFunc("/asd", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("POST")
}
