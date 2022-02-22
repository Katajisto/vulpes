package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
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

	if err != nil {
		return
	}

	type RenderData struct {
		Status           models.Status
		LastData         models.DataPoint
		LastUpdatePretty string
	}

	prettyLastUpdate := prettytime.Format(lastData.CreatedAt)

	renderData := RenderData{
		Status:           status,
		LastData:         lastData,
		LastUpdatePretty: prettyLastUpdate,
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// TODO: Implement data cleanup.
	// If we have a record at most an hour old. We discard the new data.
	limit := time.Now().Add(-6 * time.Hour)
	alarmLimit := time.Now().Add(-1 * time.Hour)

	// Even though we reject data from datapoint, we must check if temperature is too low.
	// TODO: IMPLEMENT SOMETHING NICER FOR THIS!
	if !latest.Model.CreatedAt.Before(limit) {
		for _, temp := range data.Temperatures {
			if temp.Value < 12 {
				c.AlarmsController.SendAlarm("LÄMPÖTILA ALLE 14C!")
			}
		}
	}

	if !latest.Model.CreatedAt.Before(limit) {
		w.WriteHeader(http.StatusOK)
		// TODO: Remove this later
		log.Println("Discarded data.")
		return
	}

	err = c.DataService.AddData(data.Temperatures)
	if err != nil {
		// We dont care about the error in our sensor data sender
		// so we just log it.
		log.Println("Data adding failed: ", err)
	}
	w.WriteHeader(http.StatusOK)
}

func (c *DataController) GetJSONTemps(w http.ResponseWriter, r *http.Request) {
	// For development purposes of webcomponents.
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:5000")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Some quite ugly logic so we don't have to write this on front-end.
	// TODO: Think about refactoring this.
	type Value struct {
		X string  `json:"x"`
		Y float64 `json:"y"`
	}

	type SensorData struct {
		Name   string  `json:"name"`
		Values []Value `json:"values"`
	}

	type SensorsData struct {
		Times   []string     `json:"times"`
		Sensors []SensorData `json:"sensors"`
	}

	temps, err := c.DataService.GetAllData()
	if err != nil {
		panic(err)
	}

	allSensors := SensorsData{}

	sensorMap := make(map[string][]Value)

	for _, point := range temps {
		allSensors.Times = append(allSensors.Times, point.Timestamp)
		for _, sensor := range point.TemperatureData {
			sensorMap[sensor.Sensor] = append(sensorMap[sensor.Sensor], Value{X: point.Timestamp, Y: sensor.Value})
		}
	}

	for k, v := range sensorMap {
		allSensors.Sensors = append(allSensors.Sensors, SensorData{Name: k, Values: v})
	}

	json.NewEncoder(w).Encode(allSensors)
}

func (c *DataController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", c.Get).Methods("GET")
	r.HandleFunc("/data", c.PostJSONData).Methods("POST")
}
