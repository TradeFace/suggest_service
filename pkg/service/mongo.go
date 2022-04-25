package service

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoService struct {
	Client   *mongo.Client
	cfg      *Config
	Database *mongo.Database
	Ctx      context.Context
}

func NewMongoService(cfg *Config) (*MongoService, error) {

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

	if mc == nil || mc.Client == nil {
		return false
	}

	err := mc.Client.Ping(mc.Ctx, readpref.Primary())

	return err == nil
}

func (mc *MongoService) Close() {
	defer func() {
		if err := mc.Client.Disconnect(mc.Ctx); err != nil {
			log.Default().Println(err)
		}
	}()
}

func (mc *MongoService) GetOne(coll string, query bson.M, result interface{}) error {

	// if !mc.IsConnected() {
	// 	return errors.New("mongo not connected")
	// }

	err := mc.Database.Collection(coll).FindOne(context.Background(), query).Decode(result)
	return err
}

func (mc *MongoService) GetAll(coll string, query bson.M, results interface{}) (err error) {

	// if !mc.IsConnected() {
	// 	return errors.New("mongo not connected")
	// }

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
