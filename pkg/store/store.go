package store

import (
	"github.com/tradeface/suggest_service/pkg/service"
)

type Stores struct {
	Product      *ProductStore
	Domain       *DomainStore
	User         *UserStore
	Auth         *AuthStore
	ElasticQuery *ElasticQueryStore
}

func New(service *service.Service) (*Stores, error) {
	return &Stores{
		Product:      NewProductStore(service.Elastic),
		Domain:       NewDomainStore(service.Mongo),
		User:         NewUserStore(service.Mongo),
		Auth:         NewAuthStore(),
		ElasticQuery: NewElasticQueryStore(),
	}, nil
}
