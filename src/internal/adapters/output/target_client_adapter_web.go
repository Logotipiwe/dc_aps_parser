package output

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"ports-adapters-study/src/internal/core/domain"
	"strconv"
)

var link = "https://www.avito.ru/js/1/map/items?categoryId=24&locationId=653240&correctorMode=0&page=1&map=eyJzZWFyY2hBcmVhIjp7ImxhdEJvdHRvbSI6NTkuOTAwMTc2NjU1Mjc4MDI0LCJsYXRUb3AiOjU5LjkwODIyNTA1NDMyMzkxNCwibG9uTGVmdCI6MzAuNTE3OTQ0MzQ1NjgwNjQzLCJsb25SaWdodCI6MzAuNTYxMzc0Njc0MDQ5ODF9LCJ6b29tIjoxNX0%3D&params%5B201%5D=1060&params%5B504%5D=5256&params%5B550%5D%5B0%5D=5705&params%5B550%5D%5B1%5D=5704&params%5B550%5D%5B2%5D=5703&params%5B550%5D%5B3%5D=5702&verticalCategoryId=1&rootCategoryId=4&localPriority=0&disabledFilters%5Bids%5D%5B0%5D=byTitle&disabledFilters%5Bslugs%5D%5B0%5D=bt&subscription%5Bvisible%5D=true&subscription%5BisShowSavedTooltip%5D=false&subscription%5BisErrorSaved%5D=false&subscription%5BisAuthenticated%5D=true&searchArea%5BlatBottom%5D=59.900176655278024&searchArea%5BlonLeft%5D=30.517944345680643&searchArea%5BlatTop%5D=59.908225054323914&searchArea%5BlonRight%5D=30.56137467404981&viewPort%5Bwidth%5D=1012&viewPort%5Bheight%5D=374&limit=50&countAndItemsOnly=1"

type TargetClientAdapterWeb struct {
}

func NewTargetClientWebAdapter() *TargetClientAdapterWeb {
	return &TargetClientAdapterWeb{}
}

// TODO paging results

func (k *TargetClientAdapterWeb) GetParseResult() (domain.ParseResult, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		return domain.ParseResult{}, err
	}
	log.Println("Requesting target...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return domain.ParseResult{}, err
	}
	defer resp.Body.Close()
	log.Println(resp.Status)
	var targetAnswer TargetAnswer
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.ParseResult{}, err
	}
	err = json.Unmarshal(body, &targetAnswer)
	if err != nil {
		return domain.ParseResult{}, err
	}
	log.Println("Got " + strconv.Itoa(targetAnswer.TotalCount) + " aps")
	result := domain.ParseResult{
		Items: make([]domain.ParseItem, 0),
	}
	for _, targetResultItem := range targetAnswer.Items {
		result.Items = append(result.Items, domain.ParseItem{
			ID:   targetResultItem.ID,
			Link: req.Host + targetResultItem.UrlPath,
		})
	}

	return result, nil
}

type TargetAnswer struct {
	Items      []TargetAnswerItem `json:"items"`
	TotalCount int                `json:"totalCount"`
}

type TargetAnswerItem struct {
	ID      int64  `json:"id"`
	UrlPath string `json:"urlPath"`
}
