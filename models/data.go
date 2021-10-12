package models

import (
	"time"
	"log"

	"gorm.io/gorm"
)

type Temperature struct {
	gorm.Model
	DataPointID int64
	Sensor      string  `json:"sensor"`
	Value       float64 `json:"value"`
}

type DataPoint struct {
	gorm.Model
	Timestamp       string        `json:"timestamp"`
	TemperatureData []Temperature `json:"temperatureData"`
}

type Status struct {
	gorm.Model
	Armed bool `json:"armed"`
}

type DataService struct {
	db *gorm.DB
}

func NewDataService(db *gorm.DB) *DataService {
	return &DataService{db: db}
}

func (d *DataService) GetStatus() (Status, error) {
	var status Status
	err := d.db.First(&status).Error
	return status, err
}

func (d *DataService) SetStatus(armed bool) error {
	var status Status
	err := d.db.First(&status).Error
	if err != nil {
		return err
	}

	status.Armed = armed
	return d.db.Save(&status).Error
}

// Gets 200 latest datapoints, in order to not overwork the frontend.
func (d *DataService) GetAllData() ([]DataPoint, error) {
	// Get amount of data points in db so we can skip most.
	var amount int64
	d.db.Model(&DataPoint{}).Count(&amount)
	log.Println("Found ", amount, " datapoints in db.")

	amount -= 200;
	if amount < 0 { amount = 0 }

	var data []DataPoint
	err := d.db.Preload("TemperatureData").Offset(int(amount)).Limit(200).Find(&data).Error
	return data, err
}

func (d *DataService) AddData(tempData []Temperature) error {
	data := DataPoint{
		Timestamp:       time.Now().Format(time.RFC3339),
		TemperatureData: tempData,
	}

	err := d.db.Create(&data).Error
	return err
}

// Get latest data point
func (d *DataService) GetLatestData() (DataPoint, error) {
	var data DataPoint
	err := d.db.Preload("TemperatureData").Last(&data).Error
	return data, err
}

func MigrateData(db *gorm.DB) {
	db.AutoMigrate(&Temperature{})
	db.AutoMigrate(&DataPoint{})
	db.AutoMigrate(&Status{})
	// If db has no status row, create one
	var status Status
	err := db.First(&status).Error
	if err != nil {
		db.Create(&Status{Armed: false})
	}
}
