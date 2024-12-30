package output

import (
	"fmt"
	drivenport "ports-adapters-study/src/internal/core/ports/output"
	"ports-adapters-study/src/internal/infrastructure/tg"
)

type tgNotificationClient struct {
	tgClient *tg.TgClient
}

func NewTgClientAdapter(
	tgClient *tg.TgClient,
) drivenport.NotificationClient {

	return &tgNotificationClient{
		tgClient: tgClient,
	}
}

func (t *tgNotificationClient) NotifyStartParsing(parserID int) error {
	return t.tgClient.SendMessage(fmt.Sprintf("Парсер %d запущен", parserID))
}

func (t *tgNotificationClient) NotifyChanges(diff int) error {
	var msg string
	if diff > 0 {
		msg = fmt.Sprintf("Квартир стало больше на %d", diff)
	} else {
		msg = fmt.Sprintf("Кввартир стало меньше на %d", -diff)
	}
	err := t.tgClient.SendMessage(msg)
	return err
}
