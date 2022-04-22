package store

import (
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/elastic"
	"github.com/tradeface/suggest_service/pkg/mongo"
)

type Stores struct {
	Product *Product
	Domain  *Domain
	User    *User
}

func New(dbconn *mongo.MongoClient, esconn *elastic.Elastic, cfg *conf.Config) (*Stores, error) {
	return &Stores{
		Product: NewProduct(dbconn, esconn, cfg),
		Domain:  NewDomain(dbconn),
		User:    NewUser(dbconn),
	}, nil
}
