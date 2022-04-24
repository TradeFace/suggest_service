package controller

import (
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/store"
)

type Provider struct {
	Domain  *DomainController
	Product *ProductController
	User    *UserController
}

func NewProvider(cfg *conf.Config, storeProvider *store.Provider) (*Provider, error) {

	productController := NewProductController(storeProvider)
	return &Provider{
		Domain: &DomainController{
			StoreProvider: storeProvider,
		},
		User: &UserController{
			StoreProvider: storeProvider,
			cfg:           cfg,
		},
		Product: productController,
	}, nil
}
