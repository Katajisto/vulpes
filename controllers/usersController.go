package controllers

import (
	"fmt"
	"net/http"

	"vulpes.ktj.st/views"
)

// Don't create this struct yourself. It needs to have the templates loaded
// before it can be used.
type UsersController struct {
	LoginView *views.View
}

func NewUsersController() *UsersController {
	return &UsersController{
		LoginView: views.NewView("main", "views/templates/login.tmpl"),
	}
}

func (uc *UsersController) Login(w http.ResponseWriter, r *http.Request) {
	uc.LoginView.Template.ExecuteTemplate(w, "main", nil)
}

func (uc *UsersController) LoginPost(w http.ResponseWriter, r *http.Request) {
	// Get username and password from form:
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Fprintln(w, username, " :: ", password)
}
