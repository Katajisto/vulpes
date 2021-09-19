package core

import (
	"github.com/gorilla/mux"
	"vulpes.ktj.st/controllers"
	"vulpes.ktj.st/core/security"
	"vulpes.ktj.st/models"
)

// Get this from env later or something.
const hmacSecret = "secrettt"

func registerRoutes(r *mux.Router) {
	db := models.InitDB()
	userService := models.NewUserService(db)
	hmac := security.NewHMAC(hmacSecret)

	usersController := controllers.NewUsersController(userService, hmac)
	usersController.RegisterRoutes(r)

}
