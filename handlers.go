package main

import (
	"net/http"

	"vulpes.ktj.st/controllers"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logged in check here or in a middleware, so we render the login page if not logged in
	// We can create the controller here, because it really doesn't matter. We host this on lambda,
	// So we can't optimize by initializing the controller higher up.
	usersController := controllers.NewUsersController()
	usersController.Login(w, r)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	usersController := controllers.NewUsersController()
	usersController.LoginPost(w, r)
}
