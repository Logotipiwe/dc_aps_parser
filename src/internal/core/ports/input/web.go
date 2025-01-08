package inputport

import (
	"dc-aps-parser/src/internal/core/application"
	"dc-aps-parser/src/internal/core/domain"
)

type ParserPort interface {
	HasActiveParser(chatID int64) bool
	LaunchParser(params domain.ParserParams) (*application.Parser, error)
	StopParser(chatID int64) error
}
