package model

import (
	"github.com/google/jsonapi"
)

type Product struct {
	ID        string `jsonapi:"primary,product"`
	ProductId string `jsonapi:"attr,productId"`
	// Suppliers []*string `jsonapi:"attr,suppliers"`
	Description string `jsonapi:"attr,description"`
}

func (product Product) JSONAPIMeta() *jsonapi.Meta {
	return &jsonapi.Meta{
		"detail": "extra details regarding the product",
	}
}
