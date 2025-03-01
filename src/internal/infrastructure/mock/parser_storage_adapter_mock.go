package mock

import (
	"dc-aps-parser/src/internal/core/domain"
)

type ParserStorageAdapterMock struct {
	parsers map[int64]domain.ParserData
}

func NewParserStorageAdapterMock() *ParserStorageAdapterMock {
	return &ParserStorageAdapterMock{
		parsers: make(map[int64]domain.ParserData),
	}
}

func (p *ParserStorageAdapterMock) SaveParser(parserData domain.ParserData) error {
	p.parsers[parserData.ChatID] = parserData
	return nil
}

func (p *ParserStorageAdapterMock) RemoveParser(chatID int64) error {
	delete(p.parsers, chatID)
	return nil
}

func (p *ParserStorageAdapterMock) GetParsers() ([]domain.ParserData, error) {
	var parsers []domain.ParserData
	for _, parserData := range p.parsers {
		parsers = append(parsers, parserData)
	}
	return parsers, nil
}

func (p *ParserStorageAdapterMock) SetParsers(parsers []domain.ParserData) {
	for _, parser := range parsers {
		p.parsers[parser.ChatID] = parser
	}
}
