package application

import (
	"ports-adapters-study/src/internal/core/domain"
	drivenport "ports-adapters-study/src/internal/core/ports/output"
)

type ResultService struct {
	resultDB     drivenport.ResultDB
	targetClient drivenport.TargetClient
}

func NewResultService(
	resultDB drivenport.ResultDB,
	targetClient drivenport.TargetClient,
) *ResultService {
	return &ResultService{
		resultDB, targetClient,
	}
}

func (s *ResultService) GetResult() (*domain.ParseResult, error) {
	result, err := s.targetClient.GetParseResult()
	if err != nil {
		return nil, err
	}
	history, err := s.GetResultHistory()
	if err != nil {
		return nil, err
	}
	result.ID = len(history) + 1
	err = s.resultDB.AddResult(*result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *ResultService) GetResultHistory() ([]domain.ParseResult, error) {
	results, err := s.resultDB.GetAllResults()
	return results, err
}
