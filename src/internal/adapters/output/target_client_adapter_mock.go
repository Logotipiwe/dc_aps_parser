package output

import (
	"ports-adapters-study/src/internal/core/domain"
	"sync"
)

type TargetClientAdapterMock struct {
	results           []domain.ParseResult
	currIndex         int
	wg                *sync.WaitGroup
	isWaitingForCalls bool
}

func NewTargetClientAdapterMock(results []domain.ParseResult) *TargetClientAdapterMock {
	return &TargetClientAdapterMock{
		results:           results,
		currIndex:         0,
		wg:                &sync.WaitGroup{},
		isWaitingForCalls: false,
	}
}

func (k *TargetClientAdapterMock) GetParseResult() (domain.ParseResult, error) {
	result := k.results[k.currIndex]
	k.currIndex++
	if k.currIndex == len(k.results) {
		k.currIndex = 0
	}
	if k.isWaitingForCalls {
		k.wg.Done()
	}
	return result, nil
}

func (k *TargetClientAdapterMock) WaitForCalls(i int) {
	k.isWaitingForCalls = true
	k.wg.Add(i)
	k.wg.Wait()
}
