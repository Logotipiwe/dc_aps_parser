package application

import (
	drivenport "ports-adapters-study/src/internal/core/ports/output"
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
	parser := newParser(
		len(p.parsers),
		p.ResultService,
		p.NotificationPort,
		p.resultStoragePort,
	)
	parser.init()
	p.parsers = append(p.parsers, parser)
	return parser, nil
}

func (p *ParserService) StopParser(ID int) error {
	p.parsers[ID].Stop()
	return nil
}

func (p *ParserService) StopAllParsers() {
	for _, parser := range p.parsers {
		parser.Stop()
	}
}

func (p *ParserService) StopAllParsersSync() []error {
	var errs []error
	for _, parser := range p.parsers {
		err := p.StopParser(parser.ID)
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
