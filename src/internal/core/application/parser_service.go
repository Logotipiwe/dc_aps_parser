package application

import (
	"fmt"
	"time"
)

type Parser struct {
	ID      int
	stopped bool
}

func (p *Parser) init() {
	go func() {
		for {
			fmt.Printf("Parsing %d with  ...\n", p.ID)
			time.Sleep(time.Second)
			if p.stopped {
				break
			}
		}
	}()
}

func (p *Parser) Stop() {
	fmt.Printf("Parser %d stopped\n", p.ID)
	p.stopped = true
}

func newParser(ID int) *Parser {
	return &Parser{
		ID:      ID,
		stopped: false,
	}
}

type ParserService struct {
	parsers []*Parser
}

func NewParserService() *ParserService {
	return &ParserService{}
}

func (p *ParserService) NewParser() (*Parser, error) {
	parser := newParser(len(p.parsers))
	parser.init()
	p.parsers = append(p.parsers, parser)
	return parser, nil
}

func (p *ParserService) StopParser(ID int) error {
	p.parsers[ID].Stop()
	return nil
}
