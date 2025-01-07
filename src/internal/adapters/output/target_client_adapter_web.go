package output

import (
	"dc-aps-parser/src/internal/core/domain"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

var hardcodeLink = "https://www.avito.ru/js/1/map/items?categoryId=24&locationId=653240&correctorMode=0&page=1&map=eyJzZWFyY2hBcmVhIjp7ImxhdEJvdHRvbSI6NTkuOTAwNzY4NjgzMTQyNzE1LCJsYXRUb3AiOjU5LjkxMzMzNTE1NzYxNjI5LCJsb25MZWZ0IjozMC40NTA1NjA4ODQ0NTYzNywibG9uUmlnaHQiOjMwLjUzODk2NjQ5MzU4NzIxOH0sInpvb20iOjE1fQ%3D%3D&params%5B201%5D=1060&params%5B504%5D=5256&params%5B550%5D%5B0%5D=5705&params%5B550%5D%5B1%5D=5704&params%5B550%5D%5B2%5D=5703&params%5B550%5D%5B3%5D=5702&verticalCategoryId=1&rootCategoryId=4&localPriority=0&disabledFilters%5Bids%5D%5B0%5D=byTitle&disabledFilters%5Bslugs%5D%5B0%5D=bt&subscription%5Bvisible%5D=true&subscription%5BisShowSavedTooltip%5D=false&subscription%5BisErrorSaved%5D=false&subscription%5BisAuthenticated%5D=true&searchArea%5BlatBottom%5D=59.900768683142715&searchArea%5BlonLeft%5D=30.45056088445637&searchArea%5BlatTop%5D=59.91333515761629&searchArea%5BlonRight%5D=30.538966493587218&viewPort%5Bwidth%5D=2060&viewPort%5Bheight%5D=584&limit=10&countAndItemsOnly=1"

// limit of page size is 50
var pageSize = 50

type TargetClientAdapterWeb struct {
}

func NewTargetClientWebAdapter() *TargetClientAdapterWeb {
	return &TargetClientAdapterWeb{}
}

func (k *TargetClientAdapterWeb) GetParseResult() (domain.ParseResult, error) {
	log.Println("Getting result...")
	targetUrl, err := url.Parse(hardcodeLink)
	if err != nil {
		return domain.NewParseResult(), err
	}
	setQueryParam(targetUrl, "limit", strconv.Itoa(pageSize))
	answers, err := k.doRequestPages(targetUrl, pageSize)
	if err != nil {
		return domain.NewParseResult(), err
	}
	result := domain.NewParseResult()
	for i, answer := range answers {
		for _, answerItem := range answer.Items {
			result.Items = append(result.Items, domain.NewParseItem(
				answerItem.ID, targetUrl.Host+answerItem.UrlPath))
		}
		log.Println(fmt.Sprintf("Merged %d items for %d answer", len(answer.Items), i+1))
	}
	log.Println("Total result items: ", len(result.Items))
	return result, nil
}

func (k *TargetClientAdapterWeb) doRequestPages(link *url.URL, pageSize int) ([]TargetAnswer, error) {
	answers := make([]TargetAnswer, 0)
	page := 1

	setPage(link, page)
	targetAnswer, err := k.doRequest(link)
	if err != nil {
		return nil, err
	}

	answers = append(answers, targetAnswer)
	requestsCount := int(math.Ceil(float64(targetAnswer.TotalCount)/float64(pageSize))) - 1

	for range requestsCount {
		page = page + 1
		setPage(link, page)
		targetAnswer, err = k.doRequest(link)
		if err != nil {
			return nil, err
		}

		answers = append(answers, targetAnswer)
	}

	return answers, nil
}

func (k *TargetClientAdapterWeb) fillResultItems(targetAnswer TargetAnswer, result *domain.ParseResult, parseRequestURI *url.URL) {
	for _, targetResultItem := range targetAnswer.Items {
		result.Items = append(result.Items, domain.ParseItem{
			ID:   targetResultItem.ID,
			Link: parseRequestURI.Host + targetResultItem.UrlPath,
		})
	}
}

func (k *TargetClientAdapterWeb) doRequest(link *url.URL) (TargetAnswer, error) {
	req, err := http.NewRequest("GET", link.String(), nil)
	if err != nil {
		return TargetAnswer{}, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TargetAnswer{}, err
	}
	defer resp.Body.Close()
	log.Println(resp.Status)
	var targetAnswer TargetAnswer
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TargetAnswer{}, err
	}
	err = json.Unmarshal(body, &targetAnswer)
	if err != nil {
		return TargetAnswer{}, err
	}

	return targetAnswer, nil
}

func setPage(url *url.URL, page int) {
	setQueryParam(url, "page", strconv.Itoa(page))
}

func setQueryParam(url *url.URL, key string, value string) {
	query := url.Query()
	query.Set(key, value)
	url.RawQuery = query.Encode()
}

type TargetAnswer struct {
	Items      []TargetAnswerItem `json:"items"`
	TotalCount int                `json:"totalCount"`
}

type TargetAnswerItem struct {
	ID      int64  `json:"id"`
	UrlPath string `json:"urlPath"`
}
