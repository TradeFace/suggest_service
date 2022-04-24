package store

import (
	"github.com/tradeface/suggest_service/pkg/service"
)

type Provider struct {
	Product      *ProductStore
	Domain       *DomainStore
	User         *UserStore
	Auth         *AuthStore
	ElasticQuery *ElasticQueryStore
}

type Config struct {
	JWTSalt string
	Service *service.Provider
}

func NewProvider(cfg *Config) (*Provider, error) {
	return &Provider{
		Product:      NewProductStore(cfg.Service.Elastic),
		Domain:       NewDomainStore(cfg.Service.Mongo),
		User:         NewUserStore(cfg.Service.Mongo, cfg.JWTSalt),
		Auth:         NewAuthStore(),
		ElasticQuery: NewElasticQueryStore(),
	}, nil
}
