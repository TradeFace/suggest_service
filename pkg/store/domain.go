package store

import (
	"errors"
	"fmt"

	"github.com/tradeface/suggest_service/pkg/document"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DomainStore struct {
	dbconn   MongoServiceInterface //*service.MongoService
	collName string
}

type MongoServiceInterface interface {
	GetAll(coll string, query bson.M, results interface{}) (err error)
	GetOne(coll string, query bson.M, result interface{}) error
	GetMongoId(id string) (objID primitive.ObjectID, err error)
}

func NewDomainStore(dbconn MongoServiceInterface) *DomainStore {

	return &DomainStore{
		dbconn:   dbconn,
		collName: "domain",
	}
}

func (d *DomainStore) GetWithId(id string) (results []*document.Domain, err error) {
	fmt.Println(id)
	objID, err := d.dbconn.GetMongoId(id)
	fmt.Println(objID, err)
	if err != nil {
		return nil, err
	}

	err = d.GetAll(bson.M{"_id": objID}, &results)
	fmt.Println("test----", err, results)
	for _, result := range results {
		result.Id, err = d.getStringId(result)
	}

	return results, err
}

func (d *DomainStore) GetWithHost(host string) (results []*document.Domain, err error) {

	//TODO: query aliases
	err = d.GetAll(bson.M{"host": host}, &results)
	for _, result := range results {
		result.Id, err = d.getStringId(result)
	}
	return results, err
}

func (d *DomainStore) GetOneWithHost(host string) (result *document.Domain, err error) {

	//TODO: query aliases
	err = d.GetOne(bson.M{"host": host}, &result)
	if err != nil {
		return result, err
	}
	fmt.Println(result, err)
	result.Id, err = d.getStringId(result)
	return result, err
}

func (d *DomainStore) GetOne(query bson.M, result interface{}) error {

	return d.dbconn.GetOne(d.collName, query, result)
}

func (d *DomainStore) GetAll(query bson.M, results interface{}) (err error) {
	fmt.Println("get in all", query)
	return d.dbconn.GetAll(d.collName, query, results)
}

func (u *DomainStore) getStringId(result *document.Domain) (string, error) {

	if result == nil {
		return "", errors.New("nil result")
	}
	return result.ObjectID.Hex(), nil
}
