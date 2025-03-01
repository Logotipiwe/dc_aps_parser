package mock

type PermissionMock struct {
	ChatID     int64
	AllowedNum int
}

type PermissionStorageAdapterMock struct {
	permissions map[int64]PermissionMock
}

func NewPermissionStorageAdapterMock() *PermissionStorageAdapterMock {
	return &PermissionStorageAdapterMock{
		permissions: make(map[int64]PermissionMock),
	}
}

func (p *PermissionStorageAdapterMock) GetPermittedApsNumForChat(chatID int64) (*int, error) {
	mock, has := p.permissions[chatID]
	if !has {
		return nil, nil
	}
	return &mock.AllowedNum, nil
}

func (p *PermissionStorageAdapterMock) SetPermissions(permissions []PermissionMock) {
	for _, permission := range permissions {
		p.permissions[permission.ChatID] = permission
	}
}
