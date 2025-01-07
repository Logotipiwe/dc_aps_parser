package input

import (
	"dc-aps-parser/src/internal/core/application"
	"dc-aps-parser/src/internal/infrastructure/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ParserAdapterTg struct {
	api           *tg.BotAPI
	parserService *application.ParserService
	adminService  *application.AdminService
}

func NewParserAdapterTg(
	api *tg.BotAPI,
	parserService *application.ParserService,
	adminService *application.AdminService,
) *ParserAdapterTg {
	adapterTg := &ParserAdapterTg{
		api:           api,
		parserService: parserService,
		adminService:  adminService,
	}
	adapterTg.initListening()
	return adapterTg
}

func (t *ParserAdapterTg) initListening() {
	t.api.ReceiveMessages(func(update tgbotapi.Update) error {
		text := update.Message.Text
		chatID := update.Message.Chat.ID
		if text == "/start" {
			return t.api.SendMessageInTg(chatID, "Start")
		}
		if text == "/stop" {
			if t.parserService.HasActiveParser(chatID) {
				err := t.parserService.StopParser(chatID)
				if err != nil {
					return t.api.SendMessageInTg(chatID, "Error stopping parser")
				}
				return t.api.SendMessageInTg(chatID, "Stopped")
			}
			return t.api.SendMessageInTg(chatID, "No parser found")

		}
		if text == "/help" {
			if t.adminService.IsAdmin(chatID) {
				return t.api.SendMessageInTg(chatID, "Admin help")
			} else {
				return t.api.SendMessageInTg(chatID, "Help")
			}
		}
		if text == "/info" {
			if t.adminService.IsAdmin(chatID) {
				return t.api.SendMessageInTg(chatID, "Admin info")
			}
		}
		if t.parserService.CanParse(text) {
			_, err := t.parserService.NewParser(chatID)
			if err != nil {
				return t.api.SendMessageInTg(chatID, "Error starting parser")
			}
			return t.api.SendMessageInTg(chatID, "Start parsing")
		}
		return t.sendUnknownMessage(chatID)
	})
}

func (t *ParserAdapterTg) sendUnknownMessage(chatID int64) error {
	return t.api.SendMessageInTg(chatID, "Unknown")
}
