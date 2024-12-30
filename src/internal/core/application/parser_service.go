package application

type ParserService struct {
	parsers []*Parser
	*ResultService
}

func NewParserService(
	resultService *ResultService,
) *ParserService {
	return &ParserService{
		parsers:       make([]*Parser, 0),
		ResultService: resultService,
	}
}

func (p *ParserService) NewParser() (*Parser, error) {
	parser := newParser(
		len(p.parsers),
		p.ResultService,
	)
	parser.init()
	p.parsers = append(p.parsers, parser)
	return parser, nil
}

func (p *ParserService) StopParser(ID int) error {
	p.parsers[ID].Stop()
	return nil
}
