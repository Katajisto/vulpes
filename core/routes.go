package core

import (
	"github.com/gorilla/mux"
	"vulpes.ktj.st/controllers"
	"vulpes.ktj.st/core/middleware"
	"vulpes.ktj.st/core/security"
	"vulpes.ktj.st/models"
)

// Get this from env later or something.
const hmacSecret = "secrettt"

func registerRoutes(r *mux.Router) {
	db := models.InitDB()
	userService := models.NewUserService(db)
	dataService := models.NewDataService(db)
	alarmsService := models.NewAlarmsService(db)
	hmac := security.NewHMAC(hmacSecret)

	usersController := controllers.NewUsersController(userService, hmac)
	dataController := controllers.NewDataController(dataService)
	alarmsController := controllers.NewAlarmsController(alarmsService, dataService)
	usersController.RegisterRoutes(r)

	auth := middleware.NewRequreUserMw(usersController.UserService)

	r.HandleFunc("/", auth.ApplyFn(dataController.Get))
	r.HandleFunc("/temperatures", dataController.PostJSONData).Methods("POST")
	r.HandleFunc("/temperatures", auth.ApplyFn(dataController.GetJSONTemps)).Methods("GET")
	r.HandleFunc("/toggleAlarm", auth.ApplyFn(dataController.ToggleArmed)).Methods("POST")
	r.HandleFunc("/telegram", auth.ApplyFn(alarmsController.TelegramOverview)).Methods("GET")
	r.HandleFunc("/telegram/add", auth.ApplyFn(alarmsController.TelegramAdd)).Methods("POST")
	r.HandleFunc("/telegram/{id}/delete", auth.ApplyFn(alarmsController.TelegramDel)).Methods("POST")
	r.HandleFunc("/telegram/test", auth.ApplyFn(alarmsController.AlarmTest)).Methods("POST")
	r.HandleFunc("/postEvent", alarmsController.PostEventData).Methods("POST")
}