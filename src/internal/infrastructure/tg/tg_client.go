package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
)

type TgClient struct {
	client *tgbotapi.BotAPI
}

func NewTgClient() *TgClient {
	botAPI, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	return &TgClient{
		botAPI,
	}
}

func (t *TgClient) SendMessage(text string) error {
	_, err := t.client.Send(tgbotapi.NewMessage(214583870, text))
	return err
}
