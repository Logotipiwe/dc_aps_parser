package application

import (
	"dc-aps-parser/src/internal/core/domain"
	drivenport "dc-aps-parser/src/internal/core/ports/output"
	"dc-aps-parser/src/internal/infrastructure"
	"fmt"
	"net/url"
	"strconv"
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
	return s.notificationClient.SendMessage(chatID, fmt.Sprintf(s.config.TgInitialApsCountFormat, i)) // "Найдено %d объявлений. Ищу новые..."
}

func (s *ParserNotificationService) SendNewApInfo(chatID int64, item domain.ParseItem) error {
	return s.notificationClient.SendMessage(chatID, "Новое объявление "+item.Link)
}

func (s *ParserNotificationService) SendUserStartMessage(chatID int64) error {
	return s.notificationClient.SendMessage(chatID, s.config.TgUserStartMessage)
}

func (s *ParserNotificationService) SendErrorStoppingParser(chatID int64) error {
	return s.notificationClient.SendMessage(chatID, s.config.TgErrorStoppingParserMessage) // Ошибка остановки парсера
}

func (s *ParserNotificationService) SendParserStopped(chatID int64) error {
	return s.notificationClient.SendMessage(chatID, s.config.TgParserStoppedMessage) // Обработка новых объявлений остановлена
}

func (s *ParserNotificationService) SendParserAlreadyStopped(chatID int64) error {
	return s.notificationClient.SendMessage(chatID, s.config.TgParserAlreadyStoppedMessage) // Обработка объявлений уже остановлена
}

func (s *ParserNotificationService) SendAdminHelp(chatID int64) error {
	return s.notificationClient.SendMessage(chatID, s.config.TgAdminHelpMessage)
}

func (s *ParserNotificationService) SendUserHelp(chatID int64) error {
	return s.notificationClient.SendMessage(chatID, s.config.TgUserHelpMessage)
}

func (s *ParserNotificationService) SendStoppedParserStatus(chatID int64) error {
	return s.notificationClient.SendMessage(chatID, s.config.TgStoppedParserStatusMessage) // Новые объявления не обрабатываются. Отправьте /help, чтобы узнать, как начать получать уведомления
}

func (s *ParserNotificationService) SendUnknownCommand(chatID int64) error {
	return s.notificationClient.SendMessage(chatID, s.config.TgUnknownCommandMessage) // Неизвестная команда, попробуйте другую
}

func (s *ParserNotificationService) SendAdminInfo(chatID int64, parsers []*Parser) error {
	text := "Активных парсеров " + strconv.Itoa(len(parsers))
	for _, parser := range parsers {
		link, err := url.QueryUnescape(parser.CurrentBrowserUrl)
		if err != nil {
			link = parser.ParseLink
		}
		text += fmt.Sprintf("\n%d %s, aps: %d. %s", parser.ChatID, parser.UserName, parser.CurrentApsCount, link)
	}
	return s.notificationClient.SendMessage(chatID, text)
}

func (s *ParserNotificationService) SendParserStatus(chatID int64) error {
	return s.notificationClient.SendMessage(chatID, s.config.TgActiveParserStatus) // Парсер включен и ожидает новых объявлений, они будут присланы сюда как только появятся. Чтобы остановить его, отправьте /stop
}

func (s *ParserNotificationService) SendApsNumNotAllowed(chatID int64, requestedNum int, allowedNum int) error {
	return s.notificationClient.SendMessage(chatID, fmt.Sprintf(s.config.TgApsNumNotAllowedFormat, requestedNum, allowedNum))
}

func (s *ParserNotificationService) SendErrorMessage(chatID int64) error {
	return s.notificationClient.SendMessage(chatID, s.config.TgErrorMessage)
}
