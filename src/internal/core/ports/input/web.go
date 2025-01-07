package inputport

import "dc-aps-parser/src/internal/core/application"

type ParserPort interface {
	HasActiveParser(chatID int64) bool
	NewParser(chatID int64, parseLink string) (*application.Parser, error)
	StopParser(chatID int64) error
}
