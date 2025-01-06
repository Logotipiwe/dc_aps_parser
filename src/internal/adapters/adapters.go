package adapters

import (
	. "ports-adapters-study/src/internal/adapters/output"
	"ports-adapters-study/src/internal/core/domain"
	"ports-adapters-study/src/internal/core/ports"
)

func CreateAdapters() ports.OutputPorts {
	return ports.OutputPorts{
		TargetClientPort: NewTargetClientWebAdapter(),
		NotificationPort: NewNotificationAdapterTg(),
	}
}

func CreateMockAdapters(results []domain.ParseResult) ports.OutputPorts {
	return ports.OutputPorts{
		TargetClientPort: NewTargetClientAdapterMock(results),
		NotificationPort: NewNotificationAdapterMock(),
	}
}
