package application

import (
	"ports-adapters-study/src/internal/core/domain"
	drivenport "ports-adapters-study/src/internal/core/ports/output"
)

type ResultService struct {
	resultDB drivenport.ResultDB
}

func NewResultService(resultDB drivenport.ResultDB) *ResultService {
	return &ResultService{
		resultDB,
	}
}

func (s *ResultService) GetResult() (*domain.ParseResult, error) {
	results, err := s.resultDB.GetAllResults()
	if err != nil {
		return nil, err
	}
	result := domain.ParseResult{
		ID:     len(results) + 1,
		ApsNum: len(results),
	}
	err = s.resultDB.AddResult(result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *ResultService) GetResultHistory() ([]domain.ParseResult, error) {
	results, err := s.resultDB.GetAllResults()
	return results, err
}
