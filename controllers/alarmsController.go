package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"vulpes.ktj.st/core/telegram"
	"vulpes.ktj.st/models"
	"vulpes.ktj.st/views"
)

type AlarmsController struct {
	AlarmsView      *views.View
	AlarmsService   *models.AlarmsService
	TelegramService *telegram.TelegramService
}

func NewAlarmsController(as *models.AlarmsService) *AlarmsController {
	return &AlarmsController{
		AlarmsView:      views.NewView("main", "views/templates/telegram.tmpl"),
		AlarmsService:   as,
		TelegramService: telegram.NewTelegramService(),
	}
}

func (as *AlarmsController) TelegramOverview(w http.ResponseWriter, r *http.Request) {
	targets := as.AlarmsService.GetTgTargets()
	as.AlarmsView.Render(w, targets)
}

func (as *AlarmsController) TelegramDel(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	// convert id to uint
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Fprintf(w, "Error converting id to int")
	}

	log.Println("DELETING USER NUM ", idInt)

	as.AlarmsService.DeleteTarget(uint64(idInt))

	http.Redirect(w, r, "/telegram", http.StatusFound)
}

func (as *AlarmsController) TelegramAdd(w http.ResponseWriter, r *http.Request) {
	// convert chatid to int
	chatid := r.FormValue("chat")
	chatidInt, err := strconv.Atoi(chatid)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(r.FormValue("name"))

	as.AlarmsService.AddTgTarget(int64(chatidInt), r.FormValue("name"))
	http.Redirect(w, r, "/telegram", http.StatusFound)
}

func (ac *AlarmsController) AlarmTest(w http.ResponseWriter, r *http.Request) {
	log.Println("ALARM TEST")
	targets := ac.AlarmsService.GetTgTargets()
	for _, target := range targets {
		go ac.TelegramService.SendMessage(target.ChatID, "Hälytyksen testiviesti!")
	}
	http.Redirect(w, r, "/telegram", http.StatusFound)
}
