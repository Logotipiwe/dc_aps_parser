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
	isWaiting          bool
	resultsService     *ResultService
	resultStorage      drivenport.ResultStoragePort
	notificationClient drivenport.NotificationPort
	isFirstParse       bool
	prevApsNum         int
}

func newParser(
	ID int,
	service *ResultService,
	notificationAdapter drivenport.NotificationPort,
	resultStorage drivenport.ResultStoragePort,
) *Parser {
	return &Parser{
		ID:                 ID,
		stopped:            false,
		stopWg:             &sync.WaitGroup{},
		isWaiting:          false,
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
			p.waitForStopOrTimeout(time.Second * 1)
			if p.stopped {
				break
			}
		}
		fmt.Printf("Parser %d finally stopped\n", p.ID)
	}()
}

func (p *Parser) waitForStopOrTimeout(duration time.Duration) {
	p.stopWg.Add(1)
	go func() {
		p.isWaiting = true
		time.Sleep(duration)
		if p.isWaiting {
			p.stopWg.Done()
			p.isWaiting = false
		}
	}()
	p.stopWg.Wait()
}

func (p *Parser) Stop() {
	if !p.stopped && p.isWaiting {
		p.stopWg.Done()
	}
	p.isWaiting = false
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
		fmt.Printf("Parser %d. First parse got %d aps\n", p.ID, result.ApsNum)
		p.isFirstParse = false
	} else {
		if p.prevApsNum != result.ApsNum {
			diff := result.ApsNum - p.prevApsNum
			var msg string
			if diff > 0 {
				msg = fmt.Sprintf("Квартир стало больше на %d", diff)
			} else {
				msg = fmt.Sprintf("Кввартир стало меньше на %d", -diff)
			}
			_ = p.notificationClient.SendMessage(msg)
			fmt.Printf("Parser %d. Num diff: %d. Total now: %d\n", p.ID, diff, result.ApsNum)
		} else {
			fmt.Printf("Parser %d. Nothing changed.\n", p.ID)
		}
	}

	err = p.resultStorage.AddResult(*result)
	if err != nil {
		_ = fmt.Errorf("Parser %d. Add result error: %s\n", p.ID, err)
	}
	p.prevApsNum = result.ApsNum
}
