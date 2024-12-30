package application

import (
	"fmt"
	drivenport "ports-adapters-study/src/internal/core/ports/output"
	"time"
)

type Parser struct {
	ID                 int
	stopped            bool
	resultsService     *ResultService
	notificationClient drivenport.NotificationClient
	isFirstParse       bool
	prevApsNum         int
}

func (p *Parser) init() {
	_ = p.notificationClient.NotifyStartParsing(p.ID)
	go func() {
		for {
			fmt.Printf("Parser %d. Parsing...\n", p.ID)
			p.doParse()
			time.Sleep(time.Second * 20)
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

func (p *Parser) doParse() {
	result, err := p.resultsService.GetResult()
	if err != nil {
		fmt.Printf("Parser %d. Get result error: %s\n", p.ID, err)
		return
	}
	if p.isFirstParse {
		fmt.Printf("Parser %d. First parse got %d aps\n", p.ID, result.ApsNum)
		p.isFirstParse = false
	} else {
		if p.prevApsNum != result.ApsNum {
			diff := result.ApsNum - p.prevApsNum
			_ = p.notificationClient.NotifyChanges(diff)
			fmt.Printf("Parser %d. Num diff: %d. Total now: %d\n", p.ID, diff, result.ApsNum)
		} else {
			fmt.Printf("Parser %d. Nothing changed.\n", p.ID)
		}
	}
	p.prevApsNum = result.ApsNum
}

func newParser(ID int, service *ResultService, client drivenport.NotificationClient) *Parser {
	return &Parser{
		ID:                 ID,
		stopped:            false,
		isFirstParse:       true,
		resultsService:     service,
		notificationClient: client,
	}
}
