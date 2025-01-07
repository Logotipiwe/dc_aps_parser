package outputport

import "dc-aps-parser/src/internal/core/domain"

type ResultStoragePort interface {
	AddResult(result domain.ParseResult) error
	GetAllResults() ([]domain.ParseResult, error)
}

type TargetClientPort interface {
	GetParseResult() (domain.ParseResult, error)
}

type NotificationPort interface {
	SendMessage(chatID int64, text string) error
	SendMessageWithImages(chatID int64, text string, images []string) error
}
