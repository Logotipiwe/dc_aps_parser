package drivenport

import "ports-adapters-study/src/internal/core/domain"

type ResultDB interface {
	AddResult(result domain.ParseResult) error
	GetAllResults() ([]domain.ParseResult, error)
}

type TargetClient interface {
	GetParseResult() (*domain.ParseResult, error)
}

type NotificationClient interface {
	NotifyStartParsing(parserID int) error
	NotifyChanges(diff int) error
}
