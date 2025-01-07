package application

type AdminService struct {
}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func (s *AdminService) IsAdmin(chatID int64) bool {
	// TODO
	return chatID == 214583870
}
