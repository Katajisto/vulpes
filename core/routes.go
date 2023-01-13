package core

import (
	"github.com/gorilla/mux"
	"log"
	"vulpes.ktj.st/controllers"
	"vulpes.ktj.st/core/middleware"
	"vulpes.ktj.st/core/security"
	"vulpes.ktj.st/models"
)

// Get this from env later or something.
const hmacSecret = "secrettt"

func RegisterRoutes(r *mux.Router) {
	db := models.InitDB()
	userService := models.NewUserService(db)
	dataService := models.NewDataService(db)
	alarmsService := models.NewAlarmsService(db)
	hmac := security.NewHMAC(hmacSecret)

	users, err := userService.GetAllUsers()
	if err != nil {
		panic(err)
	}

	if len(users) < 1 {
		log.Println("No user found, creating default user")
		var defaultUser models.User
		defaultUser.Username = "admin"
		defaultUser.Password = "$2a$12$mk0Bac09PY8GJipBWBJkcOqveRA1P5gLox3VpRSbKi4E/T0N2k5ra"
		userService.AddUser(&defaultUser)
	}

	usersController := controllers.NewUsersController(userService, hmac)
	alarmsController := controllers.NewAlarmsController(alarmsService, dataService)
	// This is not optimal, but we need functions from the alarms controller to send the alert if data has problem.
	dataController := controllers.NewDataController(dataService, alarmsController)

	usersController.RegisterRoutes(r)

	auth := middleware.NewRequreUserMw(usersController.UserService)

	r.HandleFunc("/", auth.ApplyFn(dataController.Get))
	r.HandleFunc("/temperatures", dataController.PostJSONData).Methods("POST")
	r.HandleFunc("/toggleAlarm", auth.ApplyFn(dataController.ToggleArmed)).Methods("POST")
	r.HandleFunc("/telegram", auth.ApplyFn(alarmsController.TelegramOverview)).Methods("GET")
	r.HandleFunc("/telegram/add", auth.ApplyFn(alarmsController.TelegramAdd)).Methods("POST")
	r.HandleFunc("/telegram/{id}/delete", auth.ApplyFn(alarmsController.TelegramDel)).Methods("POST")
	r.HandleFunc("/telegram/test", auth.ApplyFn(alarmsController.AlarmTest)).Methods("POST")
	r.HandleFunc("/postEvent", alarmsController.PostEventData).Methods("POST")
}
