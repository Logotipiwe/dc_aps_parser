package outputport

import (
	"dc-aps-parser/src/internal/core/domain"
)

type ResultStoragePort interface {
	AddResult(result domain.ParseResult) error
	GetAllResults() ([]domain.ParseResult, error)
}

type TargetClientPort interface {
	GetParseResult(parseLink string) (domain.ParseResult, error)
}

type NotificationPort interface {
	SendMessage(chatID int64, text string) error
	SendMessageWithImages(chatID int64, text string, images []string) error
}

type ParsersStoragePort interface {
	SaveParser(parser domain.ParserData) error
	RemoveParser(parser domain.ParserData) error
	GetParsers() ([]domain.ParserData, error)
}
