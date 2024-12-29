package domain

type ParseResult struct {
	ID     int
	ApsNum int
}

func NewParseResult(id int, apsNum int) ParseResult {
	return ParseResult{id, apsNum}
}
