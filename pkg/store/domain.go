package store

import (
	"github.com/tradeface/suggest_service/pkg/document"
	"github.com/tradeface/suggest_service/pkg/service"
	"go.mongodb.org/mongo-driver/bson"
)

type DomainStore struct {
	dbconn   *service.MongoService
	collName string
}

func NewDomainStore(dbconn *service.MongoService) *DomainStore {

	return &DomainStore{
		dbconn:   dbconn,
		collName: "domain",
	}
}

func (d *DomainStore) GetWithId(id string) (results []*document.Domain, err error) {

	objID, err := d.dbconn.GetMongoId(id)
	if err != nil {
		return nil, err
	}

	err = d.getAll(bson.M{"_id": objID}, &results)
	for _, result := range results {
		d.setStringId(result)
	}

	return results, err
}

func (d *DomainStore) GetWithHost(host string) (results []*document.Domain, err error) {

	//TODO: query aliases
	err = d.getAll(bson.M{"host": host}, &results)
	for _, result := range results {
		d.setStringId(result)
	}
	return results, err
}

func (d *DomainStore) GetOneWithHost(host string) (result *document.Domain, err error) {

	//TODO: query aliases
	err = d.getOne(bson.M{"host": host}, &result)
	d.setStringId(result)
	return result, err
}

func (d *DomainStore) getOne(query bson.M, result interface{}) error {

	return d.dbconn.GetOne(d.collName, query, result)
}

func (d *DomainStore) getAll(query bson.M, results interface{}) (err error) {
	return d.dbconn.GetAll(d.collName, query, results)
}

func (d *DomainStore) setStringId(result *document.Domain) {
	if result == nil {
		return
	}
	result.Id = result.ObjectID.Hex()
}
