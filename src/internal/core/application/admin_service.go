package application

import "dc-aps-parser/src/internal/infrastructure"

type AdminService struct {
	config *infrastructure.Config
}

func NewAdminService(
	config *infrastructure.Config,
) *AdminService {
	return &AdminService{
		config: config,
	}
}

func (s *AdminService) IsAdmin(chatID int64) bool {
	return chatID == s.config.TgAdminChatId
}
