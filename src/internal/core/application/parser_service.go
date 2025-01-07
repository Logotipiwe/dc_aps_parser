package application

import (
	drivenport "dc-aps-parser/src/internal/core/ports/output"
	"dc-aps-parser/src/pkg"
	"errors"
	"github.com/google/uuid"
	"strings"
	"sync"
)

type ParserService struct {
	parsers         []*Parser
	parsersByChatID map[int64]*Parser
	*ResultService
	drivenport.NotificationPort
}

func NewParserService(
	resultService *ResultService,
	notificationPort drivenport.NotificationPort,
) *ParserService {
	return &ParserService{
		parsers:          make([]*Parser, 0),
		parsersByChatID:  make(map[int64]*Parser),
		ResultService:    resultService,
		NotificationPort: notificationPort,
	}
}

func (p *ParserService) NewParser(chatID int64) (*Parser, error) {
	if p.parsersByChatID[chatID] != nil {
		return nil, errors.New("parser already exists")
	}
	wg := new(sync.WaitGroup)
	parser := newParser(
		uuid.New().String(),
		chatID,
		wg,
		p.ResultService,
		p.NotificationPort,
	)
	parser.init()
	p.parsers = append(p.parsers, parser)
	p.parsersByChatID[chatID] = parser
	return parser, nil
}

func (p *ParserService) HasActiveParser(chatID int64) bool {
	return p.parsersByChatID[chatID] != nil
}

func (p *ParserService) StopParser(chatID int64) error {
	parser := p.parsersByChatID[chatID]
	if parser == nil {
		return errors.New("parser does not exist")
	}
	p.stopParserInternally(parser)
	return nil
}

func (p *ParserService) StopAllParsersSync() {
	for _, parser := range p.parsers {
		p.stopParserInternally(parser)
	}
	for _, parser := range p.parsers {
		parser.stopWg.Wait()
	}
}

func (p *ParserService) stopParserInternally(parser *Parser) {
	parser.Stop()
	delete(p.parsersByChatID, parser.chatID)
	pkg.RemoveElement(p.parsers, parser)
}

func (p *ParserService) CanParse(url string) bool {
	return strings.HasPrefix(url, "https://www.avito.ru/js/1/map/items")
}
