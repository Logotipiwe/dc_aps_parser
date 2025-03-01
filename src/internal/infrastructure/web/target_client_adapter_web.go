package web

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

// limit of page size is 50
var pageSize = 50

type TargetClientAdapterWeb struct {
}

func NewTargetClientWebAdapter() *TargetClientAdapterWeb {
	return &TargetClientAdapterWeb{}
}

func (k *TargetClientAdapterWeb) GetTotalCount(parseLink string) (int, error) {
	link, err := k.prepareLink(parseLink)
	if err != nil {
		return 0, err
	}
	answer, _, err := k.doRequest(link)
	if err != nil {
		return 0, err
	}
	return answer.TotalCount, nil
}

func (k *TargetClientAdapterWeb) GetParseResult(parseLink string) (domain.ParseResult, error) {
	targetUrl, err := k.prepareLink(parseLink)
	if err != nil {
		return domain.ParseResult{}, err
	}
	answers, rawAnswers, err := k.doRequestPages(targetUrl, pageSize)
	if err != nil {
		return domain.ParseResult{}, err
	}
	result := domain.NewParseResult()
	for i, answer := range answers {
		result.BrowserUrl = targetUrl.Host + answer.Url
		result.TotalCount = answer.TotalCount
		for _, answerItem := range answer.Items {
			result.Items = append(result.Items, domain.NewParseItem(
				answerItem.ID, targetUrl.Host+answerItem.UrlPath))
		}
		log.Println(fmt.Sprintf("Merged %d items for %d answer", len(answer.Items), i+1))
		if len(answer.Items) == 0 {
			log.Println(fmt.Sprintf("No items found for %v", targetUrl))
		}
	}
	result.RawItemsById = k.associateRawItemsById(rawAnswers)
	log.Println("Total result items: ", len(result.Items))
	return result, nil
}

func (k *TargetClientAdapterWeb) prepareLink(parseLink string) (*url.URL, error) {
	targetUrl, err := url.Parse(parseLink)
	if err != nil {
		return nil, err
	}
	setQueryParam(targetUrl, "limit", strconv.Itoa(pageSize))
	return targetUrl, nil
}

func (k *TargetClientAdapterWeb) doRequestPages(link *url.URL, pageSize int) ([]TargetAnswer, []map[string]interface{}, error) {
	answers := make([]TargetAnswer, 0)
	rawItems := make([]map[string]interface{}, 0)
	page := 1

	setPage(link, page)
	log.Println("Requesting total count...")
	firstAnswer, _, err := k.doRequest(link)
	if err != nil {
		return nil, nil, err
	}
	requestsCount := int(math.Ceil(float64(firstAnswer.TotalCount) / float64(pageSize)))
	for range requestsCount {
		setPage(link, page)
		log.Printf("Requesting %d page...\n", page)
		targetAnswer, rawAnswer, err := k.doRequest(link)
		if err != nil {
			return nil, nil, err
		}
		page = page + 1

		answers = append(answers, targetAnswer)
		rawItems = append(rawItems, rawAnswer.Items...)
	}

	return answers, rawItems, nil
}

func (k *TargetClientAdapterWeb) fillResultItems(targetAnswer TargetAnswer, result *domain.ParseResult, parseRequestURI *url.URL) {
	for _, targetResultItem := range targetAnswer.Items {
		result.Items = append(result.Items, domain.ParseItem{
			ID:   targetResultItem.ID,
			Link: parseRequestURI.Host + targetResultItem.UrlPath,
		})
	}
}

func (k *TargetClientAdapterWeb) doRequest(link *url.URL) (TargetAnswer, *RawAnswer, error) {
	req, err := http.NewRequest("GET", link.String(), nil)
	if err != nil {
		return TargetAnswer{}, nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TargetAnswer{}, nil, err
	}
	defer resp.Body.Close()
	log.Println(resp.Status)
	var targetAnswer TargetAnswer
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TargetAnswer{}, nil, err
	}
	err = json.Unmarshal(body, &targetAnswer)
	if err != nil {
		return TargetAnswer{}, nil, err
	}

	var rawAnswer RawAnswer
	_ = json.Unmarshal(body, &rawAnswer)

	return targetAnswer, &rawAnswer, nil
}

func (k *TargetClientAdapterWeb) associateRawItemsById(rawAnswers []map[string]interface{}) map[int64]map[string]interface{} {
	result := make(map[int64]map[string]interface{})
	for _, rawItem := range rawAnswers {
		if id, ok := rawItem["id"].(float64); ok {
			result[int64(id)] = rawItem
		} else {
			log.Printf("Failed to parse id %v of result item\n", rawItem["id"])
		}
	}
	return result
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
	Url        string             `json:"url"`
	Items      []TargetAnswerItem `json:"items"`
	TotalCount int                `json:"totalCount"`
}

type TargetAnswerItem struct {
	ID      int64  `json:"id"`
	UrlPath string `json:"urlPath"`
}

type RawAnswer struct {
	Items []map[string]interface{} `json:"items"`
}
