package application

import (
	drivenport "dc-aps-parser/src/internal/core/ports/output"
	"fmt"
	"sync"
	"time"
)

type Parser struct {
	ID                 string
	chatID             int64
	parseInterval      time.Duration
	stopped            bool
	stopWg             *sync.WaitGroup
	resultsService     *ResultService
	notificationClient drivenport.NotificationPort
	isFirstParse       bool
	prevApsNum         int
}

func newParser(ID string, chatID int64, parseInterval time.Duration, stopWg *sync.WaitGroup, service *ResultService, notificationAdapter drivenport.NotificationPort) *Parser {
	stopWg.Add(1)
	return &Parser{
		ID:                 ID,
		chatID:             chatID,
		parseInterval:      parseInterval,
		stopped:            false,
		stopWg:             stopWg,
		isFirstParse:       true,
		resultsService:     service,
		notificationClient: notificationAdapter,
	}
}

func (p *Parser) init() {
	_ = p.notificationClient.SendMessage(p.chatID, "Парсер запущен")
	go func() {
		for {
			fmt.Printf("Parser %d. Parsing...\n", p.ID)
			p.doParse()
			time.Sleep(p.parseInterval)
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
		_ = p.notificationClient.SendMessage(p.chatID, fmt.Sprintf("Найдено %d объявлений. Ищу новые...", len(result.Items)))
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
			_ = p.notificationClient.SendMessage(p.chatID, msg)
			fmt.Printf("Parser %d. Num diff: %d. Total now: %d\n", p.ID, diff, len(result.Items))
		} else {
			fmt.Printf("Parser %d. Nothing changed.\n", p.ID)
		}
	}

	p.prevApsNum = len(result.Items)
}
