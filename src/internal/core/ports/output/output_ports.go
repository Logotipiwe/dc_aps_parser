package outputport

import (
	"dc-aps-parser/src/internal/core/domain"
)

type ResultStoragePort interface {
	SaveNewRawItem(apId int64, rawItem map[string]interface{}) error
}

type TargetClientPort interface {
	GetParseResult(parseLink string) (domain.ParseResult, error)
	GetTotalCount(parseLink string) (int, error)
}

type NotificationPort interface {
	SendMessage(chatID int64, text string) error
	SendMessageWithImages(chatID int64, text string, images []string) error
}

type ParsersStoragePort interface {
	SaveParser(parser domain.ParserData) error
	RemoveParser(chatID int64) error
	GetParsers() ([]domain.ParserData, error)
}

type PermissionsStoragePort interface {
	GetPermittedApsNumForChat(chatID int64) (*int, error)
}
