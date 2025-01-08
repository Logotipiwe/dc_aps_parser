package output

import (
	"dc-aps-parser/src/internal/core/domain"
	"strconv"
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
		results: []domain.ParseResult{
			{Items: make([]domain.ParseItem, 0)},
		},
	}
}

func (t *TargetClientAdapterMock) GetTotalCount(parseLink string) (int, error) {
	if len(t.results) == 0 {
		return 0, nil
	}
	return t.results[0].TotalCount, nil
}

func (t *TargetClientAdapterMock) SetResults(results []int) {
	t.results = []domain.ParseResult{}

	for i, itemsNum := range results {
		result := domain.NewParseResult()
		result.TotalCount = itemsNum
		result.BrowserUrl = "browser_url_" + strconv.Itoa(i+1)
		for j := range itemsNum {
			item := domain.NewParseItem(int64(j), "link_"+strconv.Itoa(j+1))
			result.Items = append(result.Items, item)
		}
		t.results = append(t.results, result)
	}
}

func (t *TargetClientAdapterMock) GetParseResult(string) (domain.ParseResult, error) {
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
