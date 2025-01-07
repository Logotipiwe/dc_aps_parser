package application

import (
	drivenport "dc-aps-parser/src/internal/core/ports/output"
	"strings"
	"sync"
)

type ParserService struct {
	parsers []*Parser
	*ResultService
	drivenport.NotificationPort
}

func NewParserService(
	resultService *ResultService,
	notificationPort drivenport.NotificationPort,
) *ParserService {
	return &ParserService{
		parsers:          make([]*Parser, 0),
		ResultService:    resultService,
		NotificationPort: notificationPort,
	}
}

func (p *ParserService) NewParser() (*Parser, error) {
	wg := new(sync.WaitGroup)
	parser := newParser(
		len(p.parsers),
		wg,
		p.ResultService,
		p.NotificationPort,
	)
	parser.init()
	p.parsers = append(p.parsers, parser)
	return parser, nil
}

func (p *ParserService) StopParser(ID int) {
	p.parsers[ID].Stop()
}

func (p *ParserService) StopAllParsersSync() {
	for _, parser := range p.parsers {
		parser.Stop()
	}
	for _, parser := range p.parsers {
		parser.stopWg.Wait()
	}
}

func (p *ParserService) CanParse(url string) bool {
	return strings.HasPrefix(url, "https://www.avito.ru/js/1/map/items")
}
