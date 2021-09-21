package controllers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"vulpes.ktj.st/models"
	"vulpes.ktj.st/views"
)

type DataController struct {
	DataView    *views.View
	DataService *models.DataService
}

func NewDataController(ds *models.DataService) *DataController {
	return &DataController{
		DataView:    views.NewView("main", "views/templates/data.tmpl"),
		DataService: ds,
	}
}

func (c *DataController) Get(w http.ResponseWriter, r *http.Request) {
	data, err := c.DataService.GetAllData()
	if err != nil {
		return
	}
	err = c.DataView.Render(w, data)
	if err != nil {
		log.Println(err)
	}
}

func (c *DataController) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/", c.Get).Methods("GET")
}
