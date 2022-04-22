package mongo

import (
	"context"
	"log"
	"time"

	"github.com/tradeface/suggest_service/internal/conf"
	mongo_driver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoClient struct {
	Client   *mongo_driver.Client
	cfg      *conf.Config
	Database *mongo_driver.Database
	Ctx      context.Context
}

func NewClient(cfg *conf.Config) (*MongoClient, error) {

	mc := &MongoClient{
		cfg: cfg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	mc.Ctx = ctx

	return mc, mc.Connect()
}

func (mc *MongoClient) Connect() (err error) {

	mc.Client, err = mongo_driver.Connect(
		mc.Ctx,
		options.Client().ApplyURI(mc.cfg.MongoURI),
	)

	if err == nil {
		mc.Database = mc.Client.Database(mc.cfg.MongoDB)
	}

	return err
}

func (mc *MongoClient) IsConnected() bool {

	err := mc.Client.Ping(mc.Ctx, readpref.Primary())
	if err != nil {
		return false
	}
	return true
}

func (mc *MongoClient) close() {
	defer func() {
		if err := mc.Client.Disconnect(mc.Ctx); err != nil {
			log.Default().Println(err)
		}
	}()
}
