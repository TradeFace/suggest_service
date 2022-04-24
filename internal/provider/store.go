package provider

// Not in use, demo purpose

import (
	"github.com/tradeface/suggest_service/pkg/store"
)

type StoreProvider struct {
	Product      *store.ProductStore
	Domain       *store.DomainStore
	User         *store.UserStore
	Auth         *store.AuthStore
	ElasticQuery *store.ElasticQueryStore
}

type StoreConfig struct {
	JWTSalt         string
	ServiceProvider *ServiceProvider
}

func NewStoreProvider(cfg *StoreConfig) (*StoreProvider, error) {
	return &StoreProvider{
		Product:      store.NewProductStore(cfg.ServiceProvider.Elastic),
		Domain:       store.NewDomainStore(cfg.ServiceProvider.Mongo),
		User:         store.NewUserStore(cfg.ServiceProvider.Mongo, cfg.JWTSalt),
		Auth:         store.NewAuthStore(),
		ElasticQuery: store.NewElasticQueryStore(),
	}, nil
}
