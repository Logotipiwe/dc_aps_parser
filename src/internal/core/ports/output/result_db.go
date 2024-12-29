package drivenport

import "ports-adapters-study/src/internal/core/domain"

type ResultDB interface {
	AddResult(result domain.ParseResult) error
	GetAllResults() ([]domain.ParseResult, error)
}
