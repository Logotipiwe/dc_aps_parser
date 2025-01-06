package output

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

type NotificationAdapterTg struct {
	bot *tgbotapi.BotAPI
}

func (n *NotificationAdapterTg) SendMessage(text string) error {
	_, err := n.bot.Send(tgbotapi.NewMessage(214583870, text))
	return err
}

func NewNotificationAdapterTg() *NotificationAdapterTg {
	botAPI, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	return &NotificationAdapterTg{
		botAPI,
	}
}
