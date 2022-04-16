package store

import (
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/elastic"
	"github.com/tradeface/suggest_service/pkg/mongo"
)

type Stores struct {
	Domain  DomainStore
	Product ProductStore
}

// New create all the stores
func New(mongo *mongo.MongoClient, es *elastic.ElasticClient, cfg *conf.Config) (*Stores, error) {
	return &Stores{
		Domain:  NewDomainStore(mongo, cfg),
		Product: NewProductStore(es, cfg),
	}, nil
}
