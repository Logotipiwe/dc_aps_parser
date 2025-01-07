package application

import (
	"dc-aps-parser/src/internal/core/domain"
	"fmt"
	"sync"
	"time"
)

type Parser struct {
	ID                        string
	chatID                    int64
	parseInterval             time.Duration
	parseLink                 string
	stopped                   bool
	stopWg                    *sync.WaitGroup
	resultsService            *ResultService
	parserNotificationService *ParserNotificationService
	isFirstParse              bool
	apsMemory                 map[int64]domain.ParseItem
}

func newParser(
	ID string,
	chatID int64,
	parseLink string,
	parseInterval time.Duration,
	stopWg *sync.WaitGroup,
	service *ResultService,
	parserNotificationService *ParserNotificationService,
) *Parser {
	stopWg.Add(1)
	return &Parser{
		ID:                        ID,
		chatID:                    chatID,
		parseInterval:             parseInterval,
		parseLink:                 parseLink,
		stopped:                   false,
		stopWg:                    stopWg,
		isFirstParse:              true,
		resultsService:            service,
		parserNotificationService: parserNotificationService,
		apsMemory:                 make(map[int64]domain.ParseItem),
	}
}

func (p *Parser) init() {
	_ = p.parserNotificationService.SendParserLaunched(p.chatID)
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
	result, err := p.resultsService.GetResult(p.parseLink)
	if err != nil {
		return
	}

	for _, item := range result.Items {
		if _, has := p.apsMemory[item.ID]; !has {
			p.apsMemory[item.ID] = item
			if !p.isFirstParse {
				_ = p.parserNotificationService.SendNewApInfo(p.chatID, item)
			}
		}
	}

	if p.isFirstParse {
		_ = p.parserNotificationService.SendInitialApsCount(p.chatID, len(p.apsMemory))
		p.isFirstParse = false
	}
}
