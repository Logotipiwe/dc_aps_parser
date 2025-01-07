package ports

import (
	. "dc-aps-parser/src/internal/core/ports/input"
	. "dc-aps-parser/src/internal/core/ports/output"
)

type OutputPorts struct {
	TargetClientPort
	NotificationPort
	ParsersStoragePort
}

type InputPorts struct {
	ParserPort
}
