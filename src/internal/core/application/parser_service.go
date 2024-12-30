package application

import drivenport "ports-adapters-study/src/internal/core/ports/output"

type ParserService struct {
	parsers []*Parser
	*ResultService
	drivenport.NotificationClient
}

func NewParserService(
	resultService *ResultService,
	notificationClient drivenport.NotificationClient,
) *ParserService {
	return &ParserService{
		parsers:            make([]*Parser, 0),
		ResultService:      resultService,
		NotificationClient: notificationClient,
	}
}

func (p *ParserService) NewParser() (*Parser, error) {
	parser := newParser(
		len(p.parsers),
		p.ResultService,
		p.NotificationClient,
	)
	parser.init()
	p.parsers = append(p.parsers, parser)
	return parser, nil
}

func (p *ParserService) StopParser(ID int) error {
	p.parsers[ID].Stop()
	return nil
}
