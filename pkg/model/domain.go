package model

import (
	"github.com/tradeface/suggest_service/pkg/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Domain struct {
	ObjectID           primitive.ObjectID     `bson:"_id" json:"_id"`
	Id                 string                 `jsonapi:"primary,domain"`
	Active             bool                   `jsonapi:"attr,active"`
	Catalogs           []*string              `jsonapi:"attr,catalogs"`
	Host               string                 `jsonapi:"attr,host"`
	MainClassification string                 `jsonapi:"attr,mainClassification"`
	Settings           map[string]interface{} `jsonapi:"attr,settings"`
	Modules            []string               `jsonapi:"attr,modules"`
}

func (d *Domain) ModuleIsEnabled(module string) bool {
	s := helpers.NewSet()
	s.Append(d.Modules)
	return s.Contains(module)
}

func (d *Domain) GetSetting(module string, setting string) interface{} {

	return d.Settings[module].(map[string]interface{})[setting]
}
