package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"vulpes.ktj.st/core"
)

func main() {
	// Create a new router
	r := mux.NewRouter()
	core.RegisterRoutes(r)

	// Bind the router to a port
	err := http.ListenAndServe(":1337", r)
	log.Println("Exited due to error: ", err)
}
