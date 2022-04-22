package store

import (
	"context"
	"log"

	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/model"
	mongo_client "github.com/tradeface/suggest_service/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// DomainStore provides a doamin store for mongo
type DomainStore struct {
	mongo *mongo_client.MongoClient
	cfg   *conf.Config
}

// NewDomainStore new domain store
func NewDomainStore(mc *mongo_client.MongoClient, cfg *conf.Config) DomainStore {
	return DomainStore{mongo: mc, cfg: cfg}
}

func (ds *DomainStore) GetByHost(host string) map[string]interface{} {

	tmp := map[string]interface{}{
		"data": map[string]interface{}{
			"suppliers": []string{"henk", "jan", "piet", "klaas"},
		},
	}
	ds.GetDomainByHost(host)

	return tmp
}

func (ds *DomainStore) GetDomainByHost(host string) *model.Domain {
	collection := ds.mongo.Client.Database(ds.cfg.MongoDB).Collection("domain")
	cur, err := collection.Find(context.Background(), bson.M{"domain": host})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	var result *model.Domain = new(model.Domain)
	for cur.Next(context.Background()) {
		err := cur.Decode(result)
		if err != nil {
			log.Fatal(err)
		}
		raw := cur.Current
		log.Println(raw)
	}
	return result
}
