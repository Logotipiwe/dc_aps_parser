package adapters

import (
	. "ports-adapters-study/src/internal/adapters/output"
	"ports-adapters-study/src/internal/core/domain"
	"ports-adapters-study/src/internal/core/ports"
)

func CreateAdapters() ports.OutputPorts {
	return ports.OutputPorts{
		ResultStoragePort: NewResultStorageSqlite(),
		TargetClientPort:  NewTargetClientWebAdapter(),
		NotificationPort:  NewNotificationAdapterTg(),
	}
}

func CreateMockAdapters(results []*domain.ParseResult) ports.OutputPorts {
	return ports.OutputPorts{
		ResultStoragePort: NewResultStorageMock(),
		TargetClientPort:  NewTargetClientAdapterMock(results),
		NotificationPort:  NewNotificationAdapterMock(),
	}
}