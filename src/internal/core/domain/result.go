package domain

type ParseResult struct {
	BrowserUrl string
	Items      []ParseItem
	TotalCount int
}

type ParseItem struct {
	ID   int64
	Link string
}

func NewParseItem(ID int64, link string) ParseItem {
	return ParseItem{
		ID:   ID,
		Link: link,
	}
}

func NewParseResult() ParseResult {
	return ParseResult{
		Items: make([]ParseItem, 0),
	}
}
