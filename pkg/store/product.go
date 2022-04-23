package store

import (
	"log"

	"github.com/tradeface/suggest_service/pkg/document"
	"github.com/tradeface/suggest_service/pkg/service"
)

type ProductStore struct {
	esconn *service.ElasticService
}

func NewProductStore(esconn *service.ElasticService) *ProductStore {
	return &ProductStore{
		esconn: esconn,
	}
}

func (p *ProductStore) Search(query string) (results []*document.Product, err error) {

	// fmt.Println(query)
	res, err := p.esconn.Search(query)
	if err != nil {
		log.Printf("Error getting response: %s", err)
		return results, err
	}
	for _, hit := range res.Hits.Hits {
		results = append(results, hit.Source)
	}

	return results, err
}
