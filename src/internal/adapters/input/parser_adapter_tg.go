package input

import (
	"dc-aps-parser/src/internal/core/application"
	"dc-aps-parser/src/internal/infrastructure/tg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ParserAdapterTg struct {
	api                       *tg.BotAPI
	parserService             *application.ParserService
	adminService              *application.AdminService
	parserNotificationService *application.ParserNotificationService
}

func NewParserAdapterTg(
	api *tg.BotAPI,
	parserService *application.ParserService,
	adminService *application.AdminService,
	parserNotificationService *application.ParserNotificationService,
) *ParserAdapterTg {
	adapterTg := &ParserAdapterTg{
		api:                       api,
		parserService:             parserService,
		adminService:              adminService,
		parserNotificationService: parserNotificationService,
	}
	return adapterTg
}

func (t *ParserAdapterTg) InitListening() {
	t.api.ReceiveMessages(t.HandleTgUpdate)
}

func (t *ParserAdapterTg) HandleTgUpdate(update tgbotapi.Update) error {
	text := update.Message.Text
	chatID := update.Message.Chat.ID
	if text == "/start" {
		return t.parserNotificationService.SendUserStartMessage(chatID)
	}
	if text == "/stop" {
		if t.parserService.HasActiveParser(chatID) {
			err := t.parserService.StopParser(chatID)
			if err != nil {
				return t.parserNotificationService.SendErrorStoppingParser(chatID)
			}
			return t.parserNotificationService.SendParserStopped(chatID)
		}
		return t.parserNotificationService.SendParserAlreadyStopped(chatID)

	}
	if text == "/help" {
		if t.adminService.IsAdmin(chatID) {
			return t.parserNotificationService.SendAdminHelp(chatID)
		} else {
			return t.parserNotificationService.SendUserHelp(chatID)
		}
	}
	if text == "/info" {
		if t.adminService.IsAdmin(chatID) {
			return t.parserNotificationService.SendAdminInfo(chatID, t.parserService.GetActiveParsers())
		}
	}
	if text == "/status" {
		if t.parserService.HasActiveParser(chatID) {
			return t.parserNotificationService.SendParserStatus(chatID)
		} else {
			return t.parserNotificationService.SendStoppedParserStatus(chatID)
		}
	}
	if t.parserService.CanParse(text) {
		if t.parserService.HasActiveParser(chatID) {
			err := t.parserService.StopParser(chatID)
			if err != nil {
				return t.parserNotificationService.SendErrorStoppingParser(chatID)
			}
		}
		_, err := t.parserService.NewParser(chatID, text, false)
		if err != nil {
			return err
		}
		return nil
	}
	return t.parserNotificationService.SendUnknownCommand(chatID)
}
