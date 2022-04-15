package store

import (
	"github.com/tradeface/suggest_service/internal/conf"
)

// DomainMongo provides a doamin store for mongo
type DomainMongo struct {
	mongo *mongoClient
	cfg   *conf.Config
}

// NewDomain new project store
func NewDomain(mc *mongoClient, cfg *conf.Config) DomainMongo {
	return &DomainMongo{mongo: mc, cfg: cfg}
}
