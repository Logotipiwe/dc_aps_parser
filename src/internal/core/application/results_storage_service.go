package application

import outputport "dc-aps-parser/src/internal/core/ports/output"

type ResultsStorageService struct {
	repo outputport.ResultStoragePort
}

func NewResultsStorageService(
	repo outputport.ResultStoragePort,
) *ResultsStorageService {
	return &ResultsStorageService{
		repo: repo,
	}
}

func (r *ResultsStorageService) SaveNewRawItem(apId int64, rawItem map[string]interface{}) error {
	return r.repo.SaveNewRawItem(apId, rawItem)
}
