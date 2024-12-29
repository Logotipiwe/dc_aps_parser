package storage

import (
	"ports-adapters-study/src/internal/core/domain"
)

type ResultRepository struct {
	results []domain.ParseResult
}

func NewResultRepository() *ResultRepository {
	return &ResultRepository{}
}

func (r *ResultRepository) AddResult(result domain.ParseResult) error {
	r.results = append(r.results, result)
	return nil
}

func (r *ResultRepository) GetAllResults() ([]domain.ParseResult, error) {
	return r.results, nil
}
