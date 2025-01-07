package application

import (
	drivenport "dc-aps-parser/src/internal/core/ports/output"
	"dc-aps-parser/src/internal/infrastructure"
	"fmt"
)

type ParserNotificationService struct {
	config             *infrastructure.Config
	notificationClient drivenport.NotificationPort
}

func NewParserNotificationService(config *infrastructure.Config, notificationClient drivenport.NotificationPort) *ParserNotificationService {
	return &ParserNotificationService{
		config:             config,
		notificationClient: notificationClient,
	}
}

func (s *ParserNotificationService) SendParserLaunched(chatID int64) error {
	return s.notificationClient.SendMessage(chatID, s.config.TgParserLaunchMessage)
}

func (s *ParserNotificationService) SendInitialApsCount(chatID int64, i int) error {
	return s.notificationClient.SendMessage(chatID, fmt.Sprintf("Найдено %d объявлений. Ищу новые...", i))
}

func (s *ParserNotificationService) SendApsCountChange(chatID int64, diff int) error {
	var msg string
	if diff > 0 {
		msg = fmt.Sprintf("Квартир стало больше на %d", diff)
	} else {
		msg = fmt.Sprintf("Квартир стало меньше на %d", -diff)
	}
	return s.notificationClient.SendMessage(chatID, msg)
}
