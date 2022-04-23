package service

import (
	"context"
	"log"
	"time"

	"github.com/tradeface/suggest_service/internal/conf"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoService struct {
	Client   *mongo.Client
	cfg      *conf.Config
	Database *mongo.Database
	Ctx      context.Context
}

func NewMongoService(cfg *conf.Config) (*MongoService, error) {

	mc := &MongoService{
		cfg: cfg,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	mc.Ctx = ctx

	return mc, mc.Connect()
}

func (mc *MongoService) Connect() (err error) {

	mc.Client, err = mongo.Connect(
		mc.Ctx,
		options.Client().ApplyURI(mc.cfg.MongoURI),
	)

	if err == nil {
		mc.Database = mc.Client.Database(mc.cfg.MongoDB)
	}

	return err
}

func (mc *MongoService) IsConnected() bool {

	err := mc.Client.Ping(mc.Ctx, readpref.Primary())
	if err != nil {
		return false
	}
	return true
}

func (mc *MongoService) close() {
	defer func() {
		if err := mc.Client.Disconnect(mc.Ctx); err != nil {
			log.Default().Println(err)
		}
	}()
}

func (mc *MongoService) GetOne(coll string, query bson.M, result interface{}) error {

	err := mc.Database.Collection(coll).FindOne(context.Background(), query).Decode(result)
	return err
}

func (mc *MongoService) GetAll(coll string, query bson.M, results interface{}) (err error) {

	cur, err := mc.Database.Collection(coll).Find(context.Background(), query)
	if err != nil {
		return err
	}
	defer cur.Close(context.Background())

	if err = cur.All(context.Background(), results); err != nil {
		return err
	}
	return nil
}

func (mc *MongoService) GetMongoId(id string) (objID primitive.ObjectID, err error) {
	return primitive.ObjectIDFromHex(id)
}
