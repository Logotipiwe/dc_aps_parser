package output

import (
	"fmt"
	"ports-adapters-study/src/internal/core/domain"
)

type ResultStorageMock struct {
	results []domain.ParseResult
}

func NewResultStorageMock() *ResultStorageMock {
	return &ResultStorageMock{}
}

func (r *ResultStorageMock) AddResult(result domain.ParseResult) error {
	fmt.Println("Using mock repo to save result")
	r.results = append(r.results, result)
	return nil
}

func (r *ResultStorageMock) GetAllResults() ([]domain.ParseResult, error) {
	return r.results, nil
}
