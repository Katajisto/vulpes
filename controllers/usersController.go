package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"vulpes.ktj.st/core/middleware"
	"vulpes.ktj.st/core/security"
	"vulpes.ktj.st/models"
	"vulpes.ktj.st/views"
)

// Don't create this struct yourself. It needs to have the templates loaded
// before it can be used. Use the NewUsersController function to create it.
type UsersController struct {
	LoginView   *views.View
	UsersView   *views.View
	UserService *models.UserService
	hmac        security.HMAC
}

// Creates a new UsersController. This is the only way to create a new UsersController.
func NewUsersController(us *models.UserService, hmac security.HMAC) *UsersController {
	return &UsersController{
		LoginView:   views.NewView("main", "views/templates/login.tmpl"),
		UsersView:   views.NewView("main", "views/templates/users.tmpl"),
		UserService: us,
		hmac:        hmac,
	}
}

// Renders the login view.
func (uc *UsersController) Login(w http.ResponseWriter, r *http.Request) {
	uc.LoginView.Template.ExecuteTemplate(w, "main", nil)
}

// Renders a user list
func (uc *UsersController) UsersList(w http.ResponseWriter, r *http.Request) {
	users, err := uc.UserService.GetAllUsers()
	if err != nil {
		log.Println(err)
		return
	}
	uc.UsersView.Template.ExecuteTemplate(w, "main", users)
}

// Handles a login POST request.
func (uc *UsersController) LoginPost(w http.ResponseWriter, r *http.Request) {
	// Get username and password from form:
	username := r.FormValue("username")
	user, err := uc.UserService.GetUserByUsername(username)
	if err != nil {
		log.Println(err)
	}

	if !security.Compare(r.FormValue("password"), user.Password) {
		fmt.Fprintf(w, "Login failed")
		return
	}

	// Create session for user
	var tokenHash string

	rememberToken, err := security.RememberToken()
	if err != nil {
		log.Println(err)
	}

	tokenHash = uc.hmac.Hash(rememberToken)

	uc.UserService.AddSession(&models.Session{UserID: user.ID, Device: r.UserAgent(), TokenHash: tokenHash})

	cookie := http.Cookie{
		Name:  "token",
		Value: tokenHash,
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/", http.StatusFound)
}

// Handles new user add POST requests.
func (uc *UsersController) AddUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		fmt.Fprintf(w, "Username or password cannot be empty")
		return
	}

	err := uc.UserService.AddUser(&models.User{Username: username, Password: security.Hash(password)})
	if err != nil {
		fmt.Fprintf(w, "Error adding user: ", err)
		return
	}
	
	http.Redirect(w, r, "/users", http.StatusFound)
}

// Handles a delete user POST request.
func (uc *UsersController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	// convert id to uint
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprintf(w, "Error converting id to int")
	}

	log.Println("DELETING USER NUM ", idInt)

	err = uc.UserService.DeleteUser(uint(idInt))
	if err != nil {
		fmt.Fprintf(w, "Error deleting user: ", err)
	}

	http.Redirect(w, r, "/users", http.StatusFound)
}

// Registers the routes that our controller uses.
func (uc *UsersController) RegisterRoutes(r *mux.Router) {
	priv := middleware.NewRequreUserMw(uc.UserService)

	r.HandleFunc("/login", uc.Login).Methods("GET")
	r.HandleFunc("/login", uc.LoginPost).Methods("POST")
	r.HandleFunc("/users", priv.ApplyFn(uc.UsersList)).Methods("GET")
	r.HandleFunc("/users", priv.ApplyFn(uc.AddUser)).Methods("POST")
	r.HandleFunc("/users/{id}/delete", priv.ApplyFn(uc.DeleteUser)).Methods("POST")
}
