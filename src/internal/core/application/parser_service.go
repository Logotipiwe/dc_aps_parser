package application

import (
	drivenport "ports-adapters-study/src/internal/core/ports/output"
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
		p.resultStoragePort,
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
