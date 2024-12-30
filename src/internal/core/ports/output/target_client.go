package drivenport

import "ports-adapters-study/src/internal/core/domain"

type TargetClient interface {
	GetParseResult() (*domain.ParseResult, error)
}
