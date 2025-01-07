package output

import (
	"dc-aps-parser/src/internal/core/domain"
	"sync"
)

type TargetClientAdapterMock struct {
	results           []domain.ParseResult
	currIndex         int
	wg                *sync.WaitGroup
	isWaitingForCalls bool
}

func NewTargetClientAdapterMock() *TargetClientAdapterMock {
	return &TargetClientAdapterMock{
		currIndex:         0,
		wg:                &sync.WaitGroup{},
		isWaitingForCalls: false,
	}
}

func (t *TargetClientAdapterMock) SetResults(results []domain.ParseResult) {
	t.results = results
}

func (t *TargetClientAdapterMock) GetParseResult() (domain.ParseResult, error) {
	result := t.results[t.currIndex]
	t.currIndex++
	if t.currIndex == len(t.results) {
		t.currIndex = 0
	}
	if t.isWaitingForCalls {
		t.wg.Done()
	}
	return result, nil
}

func (t *TargetClientAdapterMock) WaitForCalls(i int) {
	t.isWaitingForCalls = true
	t.wg.Add(i)
	t.wg.Wait()
}
