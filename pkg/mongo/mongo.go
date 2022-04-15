package mongo

import "github.com/tradeface/suggest_service/internal/conf"

type MongoClient struct {
	cfg *conf.Config
}

func NewClient(cfg *conf.Config) (*MongoClient, error) {
	return &MongoClient{
		cfg: cfg,
	}, nil
}
