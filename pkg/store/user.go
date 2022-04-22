package store

import (
	"github.com/tradeface/suggest_service/pkg/authorization"
	"github.com/tradeface/suggest_service/pkg/document"
	"github.com/tradeface/suggest_service/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type User struct {
	dbconn   *mongo.MongoClient
	collName string
}

func NewUser(dbconn *mongo.MongoClient) *User {
	return &User{
		dbconn:   dbconn,
		collName: "user",
	}
}

func (u *User) Login(email string, password string) (results []*document.User, err error) {

	err = u.getAll(bson.M{
		"$and": []bson.M{
			bson.M{"email": email},
			bson.M{"password": password},
		},
	}, &results)
	for _, result := range results {
		u.setStringId(result)
		authorization.MakeJwt(result)
	}

	return results, err
}

func (u *User) GetWithId(id string) (results []*document.User, err error) {

	objID, err := u.dbconn.GetMongoId(id)
	if err != nil {
		return nil, err
	}

	err = u.getAll(bson.M{"_id": objID}, &results)
	for _, result := range results {
		u.setStringId(result)
		// authorization.MakeJwt(result)
	}

	return results, err
}

func (u *User) getOne(query bson.M, result interface{}) error {

	return u.dbconn.GetOne(u.collName, query, result)
}

func (u *User) getAll(query bson.M, results interface{}) (err error) {
	return u.dbconn.GetAll(u.collName, query, results)
}

func (u *User) setStringId(result *document.User) {
	if result == nil {
		return
	}
	result.Id = result.ObjectID.Hex()
}
