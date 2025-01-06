package ports

import (
	. "ports-adapters-study/src/internal/core/ports/input"
	. "ports-adapters-study/src/internal/core/ports/output"
)

type OutputPorts struct {
	ResultStoragePort
	TargetClientPort
	NotificationPort
}

type InputPorts struct {
	ResultPort
	ParserPort
}
