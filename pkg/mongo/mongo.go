package mongo_client

import (
	"context"
	"log"
	"time"

	"github.com/tradeface/suggest_service/internal/conf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	cfg    *conf.Config
	Client *mongo.Client
}

func NewClient(cfg *conf.Config) (*MongoClient, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, err
	}
	return &MongoClient{
		cfg:    cfg,
		Client: client,
	}, nil
}

func (mc *MongoClient) GetDomainByHost(host string) {
	collection := mc.Client.Database(mc.cfg.MongoDB).Collection("domain")
	cur, err := collection.Find(context.Background(), bson.M{"domain": host})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		raw := cur.Current
		log.Println(raw)
	}
}
