package application

import (
	"dc-aps-parser/src/internal/core/domain"
	drivenport "dc-aps-parser/src/internal/core/ports/output"
)

type ResultService struct {
	targetClientPort drivenport.TargetClientPort
}

func NewResultService(
	targetClient drivenport.TargetClientPort,
) *ResultService {
	return &ResultService{
		targetClient,
	}
}

func (s *ResultService) GetResult(parseLink string) (domain.ParseResult, error) {
	return s.targetClientPort.GetParseResult(parseLink)
}
