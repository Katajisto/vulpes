package controllers

import (
	"encoding/json"
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
	DataService     *models.DataService
	TelegramService *telegram.TelegramService
}

func NewAlarmsController(as *models.AlarmsService, ds *models.DataService) *AlarmsController {
	return &AlarmsController{
		AlarmsView:      views.NewView("main", "views/templates/telegram.tmpl"),
		AlarmsService:   as,
		TelegramService: telegram.NewTelegramService(),
		DataService:     ds,
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
	log.Println(ac.TelegramService)
	ac.SendAlarm("Hälytyksen testiviesti!")
	http.Redirect(w, r, "/telegram", http.StatusFound)
}

func (ac *AlarmsController) SendAlarm(alarm string) {
	targets := ac.AlarmsService.GetTgTargets()
	for _, target := range targets {
		ac.TelegramService.SendMessage(target.ChatID, alarm)
	}
}

// Struct that represents incoming JSON event data.
type event struct {
	EventType string      `json:"eventType"`
	EventData interface{} `json:"eventData"`
}

// Handle json event post.
func (ac *AlarmsController) PostEventData(w http.ResponseWriter, r *http.Request) {
	var data event
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	curStatus, err := ac.DataService.GetStatus()
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	isArmed := curStatus.Armed

	// Handle different types of events.
	switch data.EventType {
	case "DoorOpen":
		if isArmed {
			ac.SendAlarm("HÄLYTYS: Ovi avattiin!")
		}
	}
}
