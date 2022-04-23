package store

import (
	"fmt"

	"github.com/tradeface/suggest_service/pkg/authorization"
	"github.com/tradeface/suggest_service/pkg/document"
	"github.com/tradeface/suggest_service/pkg/service"
	"go.mongodb.org/mongo-driver/bson"
)

type UserStore struct {
	dbconn   *service.MongoService
	collName string
}

func NewUserStore(dbconn *service.MongoService) *UserStore {
	return &UserStore{
		dbconn:   dbconn,
		collName: "user",
	}
}

func (u *UserStore) Login(email string, password string) (results []*document.User, err error) {

	err = u.getAll(bson.M{
		"$and": []bson.M{
			{"email": email},
			{"password": password},
		},
	}, &results)

	for _, result := range results {
		u.setStringId(result)
		token, err := authorization.NewJwtWithUser(result)
		if err != nil {
			return results, fmt.Errorf("Unauthorized")
		}
		result.Token = token
	}

	return results, err
}

func (u *UserStore) GetWithId(id string) (results []*document.User, err error) {

	objID, err := u.dbconn.GetMongoId(id)
	if err != nil {
		return nil, err
	}

	err = u.getAll(bson.M{"_id": objID}, &results)
	for _, result := range results {
		u.setStringId(result)
	}

	return results, err
}

func (u *UserStore) getOne(query bson.M, result interface{}) error {

	return u.dbconn.GetOne(u.collName, query, result)
}

func (u *UserStore) getAll(query bson.M, results interface{}) (err error) {
	return u.dbconn.GetAll(u.collName, query, results)
}

func (u *UserStore) setStringId(result *document.User) {
	if result == nil {
		return
	}
	result.Id = result.ObjectID.Hex()
}
