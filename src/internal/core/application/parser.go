package application

import (
	"fmt"
	drivenport "ports-adapters-study/src/internal/core/ports/output"
	"sync"
	"time"
)

type Parser struct {
	ID                 int
	stopped            bool
	stopWg             *sync.WaitGroup
	resultsService     *ResultService
	resultStorage      drivenport.ResultStoragePort
	notificationClient drivenport.NotificationPort
	isFirstParse       bool
	prevApsNum         int
}

func newParser(
	ID int,
	stopWg *sync.WaitGroup,
	service *ResultService,
	notificationAdapter drivenport.NotificationPort,
	resultStorage drivenport.ResultStoragePort,
) *Parser {
	stopWg.Add(1)
	return &Parser{
		ID:                 ID,
		stopped:            false,
		stopWg:             stopWg,
		isFirstParse:       true,
		resultsService:     service,
		notificationClient: notificationAdapter,
		resultStorage:      resultStorage,
	}
}

func (p *Parser) init() {
	_ = p.notificationClient.SendMessage(fmt.Sprintf("Парсер %d запущен", p.ID))
	go func() {
		for {
			fmt.Printf("Parser %d. Parsing...\n", p.ID)
			p.doParse()
			time.Sleep(20 * time.Second)
			if p.stopped {
				break
			}
		}
		fmt.Printf("Parser %d finally stopped\n", p.ID)
		p.stopWg.Done()
	}()
}

func (p *Parser) Stop() {
	p.stopped = true
	fmt.Printf("Parser %d stopped\n", p.ID)
}

func (p *Parser) doParse() {
	result, err := p.resultsService.GetResult()
	if err != nil {
		fmt.Printf("Parser %d. Get result error: %s\n", p.ID, err)
		return
	}
	if p.isFirstParse {
		fmt.Printf("Parser %d. First parse got %d aps\n", p.ID, len(result.Items))
		p.isFirstParse = false
	} else {
		if p.prevApsNum != len(result.Items) {
			diff := len(result.Items) - p.prevApsNum
			var msg string
			if diff > 0 {
				msg = fmt.Sprintf("Квартир стало больше на %d", diff)
			} else {
				msg = fmt.Sprintf("Квартир стало меньше на %d", -diff)
			}
			_ = p.notificationClient.SendMessage(msg)
			fmt.Printf("Parser %d. Num diff: %d. Total now: %d\n", p.ID, diff, len(result.Items))
		} else {
			fmt.Printf("Parser %d. Nothing changed.\n", p.ID)
		}
	}

	p.prevApsNum = len(result.Items)
}
