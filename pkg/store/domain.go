package store

import (
	"context"

	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/elastic"
	"github.com/tradeface/suggest_service/pkg/model"
	"github.com/tradeface/suggest_service/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongo_driver "go.mongodb.org/mongo-driver/mongo"
)

type Domain struct {
	dbconn     *mongo.MongoClient
	esconn     *elastic.Elastic
	cfg        *conf.Config
	collection *mongo_driver.Collection
}

func NewDomain(dbconn *mongo.MongoClient, esconn *elastic.Elastic, cfg *conf.Config) *Domain {
	return &Domain{
		dbconn:     dbconn,
		esconn:     esconn,
		cfg:        cfg,
		collection: dbconn.Database.Collection("domain"),
	}
}

func (d *Domain) GetWithId(id string) (result []*model.Domain, err error) {

	objID, err := d.getMongoId(id)
	if err != nil {
		return nil, err
	}
	return d.getResults(bson.M{"_id": objID})
}

func (d *Domain) GetWithHost(host string) (results []*model.Domain, err error) {

	//TODO: query aliases
	return d.getResults(bson.M{"host": host})
}

func (d *Domain) GetOneWithHost(host string) (result *model.Domain, err error) {

	//TODO: query aliases
	return d.getResult(bson.M{"host": host})
}

func (d *Domain) getResult(query bson.M) (result *model.Domain, err error) {

	err = d.collection.FindOne(context.Background(), query).Decode(&result)
	return result, err
}

func (d *Domain) getResults(query bson.M) (results []*model.Domain, err error) {

	cur, err := d.collection.Find(context.Background(), query)
	if err != nil {
		return results, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {

		var result *model.Domain
		err := cur.Decode(&result)
		if err != nil {
			return results, err
		}
		d.setStringId(result)
		results = append(results, result)
	}

	return results, err
}

func (d *Domain) setStringId(result *model.Domain) {
	if result == nil {
		return
	}
	result.Id = result.ObjectID.Hex()
}

func (d *Domain) getMongoId(id string) (objID primitive.ObjectID, err error) {
	return primitive.ObjectIDFromHex(id)
}
