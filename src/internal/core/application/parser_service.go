package application

import (
	"dc-aps-parser/src/internal/core/domain"
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
	permissionsService        *PermissionsService
}

func NewParserService(
	config *infrastructure.Config,
	resultService *ResultService,
	parserNotificationService *ParserNotificationService,
	parsersStorage drivenport.ParsersStoragePort,
	permissionsService *PermissionsService,
) *ParserService {
	p := &ParserService{
		config:                    config,
		parsers:                   make([]*Parser, 0),
		parsersByChatID:           make(map[int64]*Parser),
		ResultService:             resultService,
		ParserNotificationService: parserNotificationService,
		parsersStorage:            parsersStorage,
		permissionsService:        permissionsService,
	}
	p.initParsersFromStorage()
	return p
}

func (p *ParserService) LaunchParser(
	chatID int64,
	parseLink string,
	isStartedFromStorage bool,
) (*Parser, error) {
	if p.parsersByChatID[chatID] != nil {
		return nil, errors.New("parser already exists")
	}

	err := p.checkIfApsNumAllowed(chatID, parseLink)
	if err != nil {
		return nil, err
	}

	parser := newParser(
		uuid.New().String(),
		chatID,
		parseLink,
		p.config.ParseInterval,
		new(sync.WaitGroup),
		p.ResultService,
		p.ParserNotificationService,
		isStartedFromStorage,
	)
	if !isStartedFromStorage {
		err = p.parsersStorage.SaveParser(parser.toData())
		if err != nil {
			return nil, err
		}
	}
	parser.init()
	p.parsers = append(p.parsers, parser)
	p.parsersByChatID[chatID] = parser
	return parser, nil
}

func (p *ParserService) checkIfApsNumAllowed(chatID int64, parseLink string) error {
	apsNum, err := p.ResultService.GetTotalCount(parseLink)
	if err != nil {
		return err
	}
	allowedApsNum, err := p.permissionsService.GetAllowedApsNum(chatID)
	if err != nil {
		return err
	}
	if apsNum > allowedApsNum {
		return domain.NewNotAllowedError(apsNum, allowedApsNum)
	}
	return nil
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
		_, err := p.LaunchParser(parserData.ChatID, parserData.Link, true)
		if err != nil {
			log.Println("Error starting parser for chat " + strconv.FormatInt(parserData.ChatID, 10))
			log.Println(err)
		}
	}
}
