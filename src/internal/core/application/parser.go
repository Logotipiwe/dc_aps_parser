package application

import (
	"dc-aps-parser/src/internal/core/domain"
	"fmt"
	"log"
	"sync"
	"time"
)

type Parser struct {
	ID                        string
	ChatID                    int64
	parseInterval             time.Duration
	ParseLink                 string
	stopped                   bool
	stopWg                    *sync.WaitGroup
	resultsService            *ResultService
	parserNotificationService *ParserNotificationService
	isFirstParse              bool
	apsMemory                 map[int64]domain.ParseItem
	CurrentApsCount           int
	CurrentBrowserUrl         string
	isCreatedFromStorage      bool
}

func newParser(
	ID string,
	chatID int64,
	parseLink string,
	parseInterval time.Duration,
	stopWg *sync.WaitGroup,
	service *ResultService,
	parserNotificationService *ParserNotificationService,
	isCreatedFromStorage bool,
) *Parser {
	stopWg.Add(1)
	return &Parser{
		ID:                        ID,
		ChatID:                    chatID,
		parseInterval:             parseInterval,
		ParseLink:                 parseLink,
		stopped:                   false,
		stopWg:                    stopWg,
		isFirstParse:              true,
		resultsService:            service,
		parserNotificationService: parserNotificationService,
		apsMemory:                 make(map[int64]domain.ParseItem),
		isCreatedFromStorage:      isCreatedFromStorage,
	}
}

func (p *Parser) init() {
	if !p.isCreatedFromStorage {
		_ = p.parserNotificationService.SendParserLaunched(p.ChatID)
	}
	if p.isCreatedFromStorage {
		log.Printf("New parser %v created from storage\n", p.ID)
	} else {
		log.Printf("New parser %v created from user\n", p.ID)
	}
	go func() {
		for {
			fmt.Printf("Parser %v. Parsing...\n", p.ID)
			p.doParse()
			time.Sleep(p.parseInterval)
			if p.stopped {
				break
			}
		}
		fmt.Printf("Parser %v finally stopped\n", p.ID)
		p.stopWg.Done()
	}()
}

func (p *Parser) Stop() {
	p.stopped = true
	fmt.Printf("Parser %d stopped\n", p.ID)
}

func (p *Parser) doParse() {
	result, err := p.resultsService.GetResult(p.ParseLink)
	if err != nil {
		return
	}

	for _, item := range result.Items {
		if _, has := p.apsMemory[item.ID]; !has {
			p.apsMemory[item.ID] = item
			if !p.isFirstParse {
				_ = p.parserNotificationService.SendNewApInfo(p.ChatID, item)
			}
		}
	}

	if p.isFirstParse {
		if !p.isCreatedFromStorage {
			_ = p.parserNotificationService.SendInitialApsCount(p.ChatID, len(p.apsMemory))
		}
		p.isFirstParse = false
	}
	p.CurrentApsCount = len(result.Items)
	p.CurrentBrowserUrl = result.BrowserUrl
}

func (p *Parser) toData() domain.ParserData {
	return domain.ParserData{
		ChatID: p.ChatID,
		Link:   p.ParseLink,
	}
}
