package store

import (
	"github.com/tradeface/pkg/conf"
)

type Stores struct {
	Domain Domain
}

// New create all the stores
func New(mongo *mongo.mongoClient, cfg *conf.Config) (*Stores, error) {
	return &Stores{
		Domain: NewDomain(mongo, cfg),
	}, nil
}
