package krishawebclient

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type KrishaWebClient struct {
}

func NewKrishaWebClient() *KrishaWebClient {
	return &KrishaWebClient{}
}

func (s *KrishaWebClient) RequestMapData() (*MapData, error) {
	req, _ := http.NewRequest("GET", "https://krisha.kz/a/ajax-map/map/arenda/kvartiry/almaty/?zoom=13&lat=43.24246&lon=76.92117&precision=6&bounds=txwwjg%2Ctxwtzh", nil)
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	log.Println("Requesting map data.json...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Println(resp.Status)

	var result MapData

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

type MapData struct {
	//IsTooManyAdverts bool      `json:"isTooManyAdverts"`
	//ListURL          string    `json:"listUrl"`
	//MetaData         *MetaData `json:"metaData"`
	NbTotal int `json:"nbTotal"`
}

//type MetaData struct {
//	CanonicalURL string `json:"canonicalUrl"`
//	Description  string `json:"description"`
//	Header       string `json:"header"`
//	Keywords     string `json:"keywords"`
//	Title        string `json:"title"`
//}
