package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/andanhm/go-prettytime"
	"github.com/gorilla/mux"
	"vulpes.ktj.st/graph"
	"vulpes.ktj.st/models"
	"vulpes.ktj.st/views"
)

type DataController struct {
	DataView         *views.View
	DataService      *models.DataService
	AlarmsController *AlarmsController
}

func NewDataController(ds *models.DataService, ac *AlarmsController) *DataController {
	return &DataController{
		DataView:         views.NewView("main", "views/templates/data.tmpl"),
		DataService:      ds,
		AlarmsController: ac,
	}
}

func (c *DataController) Get(w http.ResponseWriter, r *http.Request) {
	status, err := c.DataService.GetStatus()
	if err != nil {
		return
	}

	lastData, err := c.DataService.GetLatestData()

	type RenderData struct {
		Status           models.Status
		LastData         models.DataPoint
		LastUpdatePretty string
		PlotDataUri      template.URL
	}

	prettyLastUpdate := prettytime.Format(lastData.CreatedAt)

	renderData := RenderData{
		Status:           status,
		LastData:         lastData,
		LastUpdatePretty: prettyLastUpdate,
		PlotDataUri:      graph.GetTemperaturePlotImageDataUrl(c.DataService.GetAllData()),
	}

	err = c.DataView.Render(w, renderData)
	if err != nil {
		log.Println(err)
	}
}

type PostData struct {
	Temperatures []models.Temperature `json:"temperatures"`
}

func (c *DataController) ToggleArmed(w http.ResponseWriter, r *http.Request) {
	status, err := c.DataService.GetStatus()
	if err != nil {
		return
	}

	err = c.DataService.SetStatus(!status.Armed)
	if err != nil {
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// Handle json data post.
func (c *DataController) PostJSONData(w http.ResponseWriter, r *http.Request) {
	var data PostData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return
	}

	// Get latest data we have.
	latest, err := c.DataService.GetLatestData()
	if err != nil {
		// If we have a record at most an hour old. We discard the new data.
		limit := time.Now().Add(-6 * time.Hour)
		alarmLimit := time.Now().Add(-20 * time.Minute)

		didAlarm := false

		if latest.Model.CreatedAt.Before(alarmLimit) {
			for _, temp := range data.Temperatures {
				if temp.Value < 12 {
					c.AlarmsController.SendAlarm("LÄMPÖTILA ALLE 12C!")
					didAlarm = true
					break
				}
			}
		}

		// Dont discard data if there was alert.
		if !didAlarm && !latest.Model.CreatedAt.Before(limit) {
			w.WriteHeader(http.StatusOK)
			// TODO: Remove this later
			log.Println("Discarded data.")
			return
		}
	}

	err = c.DataService.AddData(data.Temperatures)
	if err != nil {
		// We dont care about the error in our sensor data sender
		// so we just log it.
		log.Println("Data adding failed: ", err)
	}
	w.WriteHeader(http.StatusOK)
}

func (c *DataController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", c.Get).Methods("GET")
	r.HandleFunc("/data", c.PostJSONData).Methods("POST")
}
