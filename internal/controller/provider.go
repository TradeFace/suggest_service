package controller

import (
	"github.com/tradeface/suggest_service/internal/provider"
)

type Provider struct {
	Domain  *DomainController
	Product *ProductController
	User    *UserController
}

func NewProvider(storeProvider *provider.StoreProvider) (*Provider, error) {

	productController := NewProductController(storeProvider)
	return &Provider{
		Domain: &DomainController{
			StoreProvider: storeProvider,
		},
		User: &UserController{
			StoreProvider: storeProvider,
		},
		Product: productController,
	}, nil
}
