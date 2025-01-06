package domain

type ParseResult struct {
	Items []ParseItem
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

func NewParseResult(items []ParseItem) ParseResult {
	return ParseResult{
		Items: items,
	}
}
