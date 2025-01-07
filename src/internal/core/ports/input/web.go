package inputport

import "dc-aps-parser/src/internal/core/application"

type ParserPort interface {
	NewParser() (*application.Parser, error)
	StopParser(ID int)
}
