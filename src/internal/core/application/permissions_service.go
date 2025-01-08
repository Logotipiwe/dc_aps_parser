package application

import (
	outputport "dc-aps-parser/src/internal/core/ports/output"
	"dc-aps-parser/src/internal/infrastructure"
)

type PermissionsService struct {
	config  *infrastructure.Config
	storage outputport.PermissionsStoragePort
}

func NewPermissionsService(
	config *infrastructure.Config,
	storage outputport.PermissionsStoragePort,
) *PermissionsService {
	return &PermissionsService{
		config,
		storage,
	}
}

func (p *PermissionsService) IsApsNumAllowed(chatID int64, apsNum int) (bool, error) {
	permitted, err := p.GetAllowedApsNum(chatID)
	if err != nil {
		return false, err
	}
	return apsNum < permitted, nil
}

func (p *PermissionsService) GetAllowedApsNum(chatID int64) (int, error) {
	permittedFromStorage, err := p.storage.GetPermittedApsNumForChat(chatID)
	if err != nil {
		return 0, err
	}
	var permitted int
	if permittedFromStorage != nil {
		permitted = *permittedFromStorage
	} else {
		permitted = p.config.DefaultAllowedApsNum
	}
	return permitted, nil
}
