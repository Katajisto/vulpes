package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var homeView view

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Render the template
	log.Println("Rendering home template")
	err := homeView.Template.ExecuteTemplate(w, homeView.Layout, nil)
	if err != nil {
		log.Println("Error rendering template: ", err)
	}
	return
}

func main() {
	homeView = newView("main", "templates/home.tmpl")
	// Create a new router
	r := mux.NewRouter()
	// Add route that loads a template
	r.HandleFunc("/", homeHandler)
	// Bind the router to a port
	http.ListenAndServe(":8080", r)
}
