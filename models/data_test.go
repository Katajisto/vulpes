package models

import (
	"log"
	"os"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func CreateTestDataService() (*DataService, error) {
	err := os.Remove("test.db")
	if err != nil {
		log.Println(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
	  Logger: logger.Default.LogMode(logger.Silent),
	})

	MigrateData(db)

	if err != nil {
		return nil, err
	}

	return NewDataService(db), nil
}


func TestStatusGet(t *testing.T) {
	ds, err := CreateTestDataService()

	if err != nil {
		t.Error(err)
	}

	// Get status without status object already inserted should insert one.
	curStatus, err := ds.GetStatus()
	if err != nil {
		t.Errorf("Expected get status to give no error. Got one.")
	}

	// Check that auto-inserted status object has default values.
	if curStatus.Armed {
		t.Errorf("Default status was armed, expected unarmed.")
	}

	ds.SetStatus(true)

	curStatus, err = ds.GetStatus()
	if err != nil {
		t.Errorf("Expected get status to give no error. Got one.")
	}

	// Check that auto-inserted status object has default values.
	if !curStatus.Armed {
		t.Errorf("Default status was unarmed, expected armed.")
	}
}

func TestLatestData(t *testing.T) {
	
}