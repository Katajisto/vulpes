package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"vulpes.ktj.st/core/security"
	"vulpes.ktj.st/models"
	"vulpes.ktj.st/views"
)

// Don't create this struct yourself. It needs to have the templates loaded
// before it can be used.
type UsersController struct {
	LoginView *views.View
	UsersView *views.View
	dbHandler *gorm.DB
	hmac      security.HMAC
}

func NewUsersController(db *gorm.DB, hmac security.HMAC) *UsersController {
	return &UsersController{
		LoginView: views.NewView("main", "views/templates/login.tmpl"),
		UsersView: views.NewView("main", "views/templates/users.tmpl"),
		dbHandler: db,
		hmac:      hmac,
	}
}

func (uc *UsersController) Login(w http.ResponseWriter, r *http.Request) {
	uc.LoginView.Template.ExecuteTemplate(w, "main", nil)
}

func (uc *UsersController) UsersList(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers(uc.dbHandler)
	if err != nil {
		log.Println(err)
		return
	}
	uc.UsersView.Template.ExecuteTemplate(w, "main", users)
}

func (uc *UsersController) LoginPost(w http.ResponseWriter, r *http.Request) {
	// Get username and password from form:
	username := r.FormValue("username")
	user, err := models.GetUserByUsername(uc.dbHandler, username)
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

	models.AddSession(uc.dbHandler, &models.Session{UserID: user.ID, Device: r.UserAgent(), TokenHash: tokenHash})

	cookie := http.Cookie{
		Name:  "token",
		Value: tokenHash,
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/users", http.StatusFound)
}

func (uc *UsersController) AddUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	if username == "" || password == "" {
		fmt.Fprintf(w, "Username or password cannot be empty")
		return
	}
	models.AddUser(uc.dbHandler, &models.User{Username: username, Password: security.Hash(password)})
	http.Redirect(w, r, "/users", http.StatusFound)
}

func (uc *UsersController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	// convert id to uint
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprintf(w, "Error converting id to int")
	}

	log.Println("DELETING USER NUM ", idInt)

	err = models.DeleteUser(uc.dbHandler, uint(idInt))
	if err != nil {
		fmt.Fprintf(w, "Error deleting user: ", err)
	}

	http.Redirect(w, r, "/users", http.StatusFound)
}

func (uc *UsersController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uc.Login(w, r)
	}).Methods("GET")
	r.HandleFunc("/login", uc.LoginPost).Methods("POST")
	r.HandleFunc("/users", uc.UsersList).Methods("GET")
	r.HandleFunc("/users", uc.AddUser).Methods("POST")
	r.HandleFunc("/users/{id}/delete", uc.DeleteUser).Methods("POST")
}
