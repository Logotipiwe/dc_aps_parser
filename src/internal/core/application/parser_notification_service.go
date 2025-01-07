package application

import (
	"dc-aps-parser/src/internal/core/domain"
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
	message := s.config.TgParserLaunchMessage
	return s.notificationClient.SendMessage(chatID, message)
}

func (s *ParserNotificationService) SendInitialApsCount(chatID int64, i int) error {
	return s.notificationClient.SendMessage(chatID, fmt.Sprintf("Найдено %d объявлений. Ищу новые...", i))
}

func (s *ParserNotificationService) SendNewApInfo(chatID int64, item domain.ParseItem) error {
	return s.notificationClient.SendMessage(chatID, "Новое объявление "+item.Link)
}
