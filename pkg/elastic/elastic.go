package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/model"
)

type ElasticClient struct {
	cfg    *conf.Config
	client *elasticsearch.Client
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
	Index  string         `json:"_index"`
	Type   string         `json:"_type"`
	Id     string         `json:"_id"`
	Score  float64        `json:"_score"`
	Source *model.Product `json:"_source"`
}

func NewClient(cfg *conf.Config) (*ElasticClient, error) {

	esCfg := elasticsearch.Config{
		Addresses: []string{
			cfg.ElasticURI,
		},
		// Username: "foo",
		// Password: "bar",
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			// TLSClientConfig: &tls.Config{
			// 	MinVersion: tls.VersionTLS12,
			// },
		},
	}
	client, err := elasticsearch.NewClient(esCfg)

	return &ElasticClient{
		cfg:    cfg,
		client: client,
	}, err
}

func (es *ElasticClient) Search(id string, sourceModel interface{}) ElasticResult {

	// rih := []ElasticResultInnerHits{
	// 	Source: &model.Product{},
	// }

	// r := ElasticResult{
	// 	Hits: ElasticResultOuterHits{
	// 		Total: struct {
	// 			Value    int    "json:\"value\""
	// 			Relation string "json:\"relation\""
	// 		}{},
	// 		MaxScore: 0,
	// 		Hits:     make(ElasticResultInnerHits,0),
	// 	},
	// }
	// x := ElasticResultInnerHitsSlice{
	// 	ElasticResultInnerHits{
	// 		Source: &sourceModel,
	// 	},
	// }
	// hh := make(ElasticResultInnerHits, 8)
	r := ElasticResult{
		// Took:     0,
		// TimedOut: false,
		// Shards: struct {
		// 	Total      int "json:\"total\""
		// 	Successful int "json:\"successful\""
		// 	Skipped    int "json:\"skipped\""
		// 	Failed     int "json:\"failed\""
		// }{},
		// Hits: ElasticResultOuterHits{
		// 	Total: struct {
		// 		Value    int    "json:\"value\""
		// 		Relation string "json:\"relation\""
		// 	}{},
		// 	MaxScore: 0,
		// 	Hits: ElasticResultInnerHitsSlice{
		// 		ElasticResultInnerHits{
		// 			Index:  id,
		// 			Type:   "",
		// 			Id:     id,
		// 			Score:  0,
		// 			Source: &model.Product{},
		// 		},
		// 	},
		// },
	}
	var buf bytes.Buffer

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"description": id,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	res, err := es.client.Search(
		es.client.Search.WithContext(context.Background()),
		es.client.Search.WithIndex(es.cfg.ElasticDB),
		es.client.Search.WithBody(&buf),
		es.client.Search.WithTrackTotalHits(true),
		es.client.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}
	log.Println("------------")
	log.Println(res)
	log.Println("------------")
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	return r
}

func (es *ElasticClient) get(id string) {

}

func (es *ElasticClient) Post(record interface{}) {

	var wg sync.WaitGroup

	for i, title := range []string{"Test One", "Test Two", "Test drie", "Test vier", "Test 5", "Test 6", "Test 7", "Test 8"} {
		wg.Add(1)

		go func(i int, title string) {
			defer wg.Done()

			doc := `{
				"productId": "%d",
				"description": "%s",
				"images": [
					{
						"name": "image%d.png"
					},
					{
						"name": "image%d.png"
					}
				]
			}`
			doc = fmt.Sprintf(doc, i+5, title, i+5, i+2)

			// Set up the request object.
			req := esapi.IndexRequest{
				Index:      es.cfg.ElasticDB,
				DocumentID: strconv.Itoa(i + 1),
				Body:       strings.NewReader(doc),
				Refresh:    "true",
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es.client)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}(i, title)
	}
	wg.Wait()

}

func (es *ElasticClient) put(record interface{}) {

}

func (es *ElasticClient) delete(id string) {

}
