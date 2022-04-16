package store

import (
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/mongo"
)

// DomainStore provides a doamin store for mongo
type DomainStore struct {
	mongo *mongo.MongoClient
	cfg   *conf.Config
}

// NewDomainStore new domain store
func NewDomainStore(mc *mongo.MongoClient, cfg *conf.Config) DomainStore {
	return DomainStore{mongo: mc, cfg: cfg}
}

func (ds *DomainStore) GetByHost(host string) map[string]interface{} {

	tmp := map[string]interface{}{
		"data": map[string]interface{}{
			"suppliers": []string{"henk", "jan", "piet", "klaas"},
		},
	}
	return tmp
}
