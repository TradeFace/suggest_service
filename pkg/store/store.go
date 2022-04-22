package store

import (
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/elastic"
	mongo_client "github.com/tradeface/suggest_service/pkg/mongo"
)

type Stores struct {
	Domain  DomainStore
	Product ProductStore
}

// New create all the stores
func New(mongo *mongo_client.MongoClient, es *elastic.ElasticClient, cfg *conf.Config) (*Stores, error) {
	return &Stores{
		Domain:  NewDomainStore(mongo, cfg),
		Product: NewProductStore(es, cfg),
	}, nil
}
