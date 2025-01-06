package outputport

import "ports-adapters-study/src/internal/core/domain"

type ResultStoragePort interface {
	AddResult(result domain.ParseResult) error
	GetAllResults() ([]domain.ParseResult, error)
}

type TargetClientPort interface {
	GetParseResult() (*domain.ParseResult, error)
}

type NotificationPort interface {
	SendMessage(text string) error
}
