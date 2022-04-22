package store

import (
	"log"

	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/elastic"
	"github.com/tradeface/suggest_service/pkg/model"
	"github.com/tradeface/suggest_service/pkg/mongo"
)

type Product struct {
	dbconn *mongo.MongoClient
	esconn *elastic.Elastic
	cfg    *conf.Config
}

func NewProduct(dbconn *mongo.MongoClient, esconn *elastic.Elastic, cfg *conf.Config) *Product {
	return &Product{
		dbconn: dbconn,
		esconn: esconn,
		cfg:    cfg,
	}
}

func (p *Product) Search(query string) (results []*model.Product, err error) {

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
