package telegram

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramService struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramService() *TelegramService {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_KEY"))
	if err != nil {
		log.Println(err)
		return nil
	}
	return &TelegramService{bot: bot}
}

func (ts *TelegramService) SendMessage(chatID int64, text string) {

	msg := tgbotapi.NewMessage(chatID, text)

	ts.bot.Send(msg)

}
