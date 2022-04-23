package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/rs/zerolog/log"
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/document"
)

type ElasticService struct {
	Client    *elasticsearch.Client
	URI       []string
	IndexName string
	User      string
	Password  string
	debug     bool
}

type ElasticResult struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits ElasticResultOuterHits `json:"hits"`
}

type ElasticResultOuterHits struct {
	Total struct {
		Value    int    `json:"value"`
		Relation string `json:"relation"`
	} `json:"total"`
	MaxScore float64                     `json:"max_score"`
	Hits     ElasticResultInnerHitsSlice `json:"hits"`
}

type ElasticResultInnerHitsSlice []ElasticResultInnerHits

type ElasticResultInnerHits struct {
	Index  string            `json:"_index"`
	Type   string            `json:"_type"`
	Id     string            `json:"_id"`
	Score  float64           `json:"_score"`
	Source *document.Product `json:"_source"`
}

func NewElasticService(cfg *conf.Config) (*ElasticService, error) {

	es := &ElasticService{
		URI:       []string{cfg.ElasticURI},
		IndexName: cfg.ElasticIndex,
		debug:     false,
	}
	err := es.startClient()
	return es, err
}

func (es *ElasticService) startClient() error {

	var err error

	es.Client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses:     es.URI,
		RetryOnStatus: []int{502, 503, 504, 429},
		RetryBackoff:  func(i int) time.Duration { return time.Duration(i) * 100 * time.Millisecond },
		MaxRetries:    5,
		Logger:        es,
	})

	if err != nil {
		log.Printf("Error creating the client: %s", err)
	}
	return err
}

func (es *ElasticService) Search(query string) (r ElasticResult, err error) {

	res, err := es.Client.Search(
		es.Client.Search.WithContext(context.Background()),
		es.Client.Search.WithIndex(es.IndexName),
		es.Client.Search.WithBody(strings.NewReader(query)),
		es.Client.Search.WithTrackTotalHits(true),
		es.Client.Search.WithPretty(),
	)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return r, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return r, fmt.Errorf("error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			return r, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return r, fmt.Errorf("error parsing the response body: %s", err)
	}
	return r, err
}

//debug functions
func (es *ElasticService) LogRoundTrip(req *http.Request, res *http.Response, err error, ts time.Time, td time.Duration) error {
	if !es.debug {
		return nil
	}
	log.Printf("req", req)
	log.Printf("res", res)

	return nil
}

func (es *ElasticService) RequestBodyEnabled() bool {
	return es.debug
}

func (es *ElasticService) ResponseBodyEnabled() bool {
	return es.debug
}
