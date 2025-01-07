package application

import (
	drivenport "dc-aps-parser/src/internal/core/ports/output"
	"dc-aps-parser/src/internal/infrastructure"
	"dc-aps-parser/src/pkg"
	"errors"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
	"sync"
)

type ParserService struct {
	config          *infrastructure.Config
	parsers         []*Parser
	parsersByChatID map[int64]*Parser
	*ResultService
	ParserNotificationService *ParserNotificationService
	parsersStorage            drivenport.ParsersStoragePort
}

func NewParserService(
	config *infrastructure.Config,
	resultService *ResultService,
	parserNotificationService *ParserNotificationService,
	parsersStorage drivenport.ParsersStoragePort,
) *ParserService {
	p := &ParserService{
		config:                    config,
		parsers:                   make([]*Parser, 0),
		parsersByChatID:           make(map[int64]*Parser),
		ResultService:             resultService,
		ParserNotificationService: parserNotificationService,
		parsersStorage:            parsersStorage,
	}
	p.initParsersFromStorage()
	return p
}

func (p *ParserService) NewParser(chatID int64, parseLink string, isSilentStart bool) (*Parser, error) {
	if p.parsersByChatID[chatID] != nil {
		return nil, errors.New("parser already exists")
	}
	wg := new(sync.WaitGroup)
	parser := newParser(
		uuid.New().String(),
		chatID,
		parseLink,
		p.config.ParseInterval,
		wg,
		p.ResultService,
		p.ParserNotificationService,
		isSilentStart,
	)
	err := p.parsersStorage.SaveParser(parser.toData())
	if err != nil {
		return nil, err
	}
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
	err := p.parsersStorage.RemoveParser(parser.toData())
	if err != nil {
		return err
	}
	p.stopParserInternally(parser)
	return nil
}

func (p *ParserService) StopAllParsersSync() {
	activeParsers := p.parsers
	for _, parser := range activeParsers {
		p.stopParserInternally(parser)
	}
	for _, parser := range activeParsers {
		parser.stopWg.Wait()
	}
}

func (p *ParserService) stopParserInternally(parser *Parser) {
	parser.Stop()
	delete(p.parsersByChatID, parser.ChatID)
	p.parsers = pkg.RemoveElement(p.parsers, parser)
}

func (p *ParserService) CanParse(url string) bool {
	return strings.HasPrefix(url, "https://www.avito.ru/js/1/map/items")
}

func (p *ParserService) GetActiveParsers() []*Parser {
	return p.parsers
}

func (p *ParserService) initParsersFromStorage() {
	parsersData, err := p.parsersStorage.GetParsers()
	if err != nil {
		log.Fatal(err)
	}
	for _, parserData := range parsersData {
		_, err := p.NewParser(parserData.ChatID, parserData.Link, true)
		if err != nil {
			log.Println("Error starting parser for chat " + strconv.FormatInt(parserData.ChatID, 10))
			log.Println(err)
		}
	}
}
