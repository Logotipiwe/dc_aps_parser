package application

import (
	. "ports-adapters-study/src/internal/core/ports"
)

type App struct {
	*ResultService
	*ParserService
}

func NewApp(a OutputPorts) App {
	resultService := NewResultService(a.ResultStoragePort, a.TargetClientPort)

	return App{
		ResultService: resultService,
		ParserService: NewParserService(resultService, a.NotificationPort),
	}
}
