package model

import (
	"github.com/google/jsonapi"
)

type Domain struct {
	ID        string    `jsonapi:"primary,domain"`
	Suppliers []*string `jsonapi:"attr,suppliers"`
}

func (domain Domain) JSONAPIMeta() *jsonapi.Meta {
	return &jsonapi.Meta{
		"detail": "extra details regarding the domain",
	}
}
