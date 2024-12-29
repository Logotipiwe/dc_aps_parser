package storage

import (
	"ports-adapters-study/src/internal/core/domain"
)

type ResultRepositoryMock struct {
}

func NewResultRepositoryMock() *ResultRepositoryMock {
	return &ResultRepositoryMock{}
}

func (r ResultRepositoryMock) AddResult(result domain.ParseResult) error {
	//TODO implement me
	return nil
}

func (r ResultRepositoryMock) GetAllResults() ([]domain.ParseResult, error) {
	//TODO implement me
	return make([]domain.ParseResult, 0), nil
}
