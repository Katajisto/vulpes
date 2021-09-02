package main

import (
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Set request headers
	w.Header().Set("Content-Type", "text/html")
	// Render the template
	log.Println("Rendering home template")
	err := homeView.Template.ExecuteTemplate(w, homeView.Layout, nil)
	if err != nil {
		log.Println("Error rendering template: ", err)
	}
}
