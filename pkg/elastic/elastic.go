package elastic

import "github.com/tradeface/suggest_service/internal/conf"

type ElasticClient struct {
	cfg *conf.Config
}

func NewClient(cfg *conf.Config) (*ElasticClient, error) {
	return &ElasticClient{
		cfg: cfg,
	}, nil
}
