package models

import (
	"time"

	"gorm.io/gorm"
)

type Temperature struct {
	gorm.Model
	DataPointID int64
	Sensor      string
	Value       float64
}

type DataPoint struct {
	gorm.Model
	Timestamp       string
	TemperatureData []Temperature
}

type DataService struct {
	db *gorm.DB
}

func NewDataService(db *gorm.DB) *DataService {
	return &DataService{db: db}
}

func (d *DataService) GetAllData() ([]DataPoint, error) {
	var data []DataPoint
	err := d.db.Find(&data).Error
	return data, err
}

func (d *DataService) AddData() error {
	data := DataPoint{
		Timestamp: time.Now().String(),
		TemperatureData: []Temperature{
			{
				Sensor: "sensor1",
				Value:  17.2,
			},
			{
				Sensor: "sensor2",
				Value:  20.0,
			},
		},
	}
	err := d.db.Create(&data).Error
	return err
}

func MigrateData(db *gorm.DB) {
	db.AutoMigrate(&Temperature{})
	db.AutoMigrate(&DataPoint{})
}
