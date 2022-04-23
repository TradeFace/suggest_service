package store

import (
	"github.com/tradeface/suggest_service/pkg/elastic"
	"github.com/tradeface/suggest_service/pkg/mongo"
)

type Stores struct {
	Product *ProductStore
	Domain  *DomainStore
	User    *UserStore
	Auth    *AuthStore
}

func New(dbconn *mongo.MongoClient, esconn *elastic.Elastic) (*Stores, error) {
	return &Stores{
		Product: NewProductStore(esconn),
		Domain:  NewDomainStore(dbconn),
		User:    NewUserStore(dbconn),
		Auth:    NewAuthStore(),
	}, nil
}
