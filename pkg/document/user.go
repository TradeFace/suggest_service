package document

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ObjectID  primitive.ObjectID `bson:"_id" json:"_id"`
	Id        string             `jsonapi:"primary,user"`
	Name      string             `jsonapi:"attr,name"`
	Email     string             `jsonapi:"attr,email"`
	Password  string             //`jsonapi:"attr,password"`
	CompanyId string             `jsonapi:"attr,companyid"`
	Roles     []string           `jsonapi:"attr,roles"`
	Token     string             `jsonapi:"attr,token,omitempty"`
}
