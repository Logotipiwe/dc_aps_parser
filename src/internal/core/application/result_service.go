package application

import (
	"ports-adapters-study/src/internal/core/domain"
	drivenport "ports-adapters-study/src/internal/core/ports/output"
)

type ResultService struct {
	resultStoragePort drivenport.ResultStoragePort
	targetClientPort  drivenport.TargetClientPort
}

func NewResultService(
	resultStorage drivenport.ResultStoragePort,
	targetClient drivenport.TargetClientPort,
) *ResultService {
	return &ResultService{
		resultStorage, targetClient,
	}
}

func (s *ResultService) GetResult() (domain.ParseResult, error) {
	return s.targetClientPort.GetParseResult()
}
