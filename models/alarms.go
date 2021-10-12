package models

import (
	"log"

	"gorm.io/gorm"
)

type AlarmsService struct {
	db *gorm.DB
}

func NewAlarmsService(db *gorm.DB) *AlarmsService {
	return &AlarmsService{db}
}

// Telegram chat to notify on alert.
type TgNotifyTarget struct {
	gorm.Model
	ChatID int64
	Name   string
}

// Gets all tg targets
func (as *AlarmsService) GetTgTargets() []TgNotifyTarget {
	var tgTargets []TgNotifyTarget
	err := as.db.Model(&TgNotifyTarget{}).Find(&tgTargets).Error
	if err != nil {
		log.Println(err)
	}

	log.Println(tgTargets)

	return tgTargets
}

// Adds a new tg target, with chatID and name
func (as *AlarmsService) AddTgTarget(chatID int64, name string) {
	log.Println(name)
	err := as.db.Create(&TgNotifyTarget{ChatID: chatID, Name: name}).Error
	if err != nil {
		log.Println(err)
	}
}

// Removes a tg target
func (as *AlarmsService) DeleteTarget(id uint64) {
	err := as.db.Delete(&TgNotifyTarget{}, id).Error
	if err != nil {
		log.Println(err)
	}
}
