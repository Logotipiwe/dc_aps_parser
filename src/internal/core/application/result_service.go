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
	return result, err
}

func (s *ResultService) GetResultHistory() ([]domain.ParseResult, error) {
	results, err := s.resultDB.GetAllResults()
	return results, err
}
