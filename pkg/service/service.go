package service

import (
	"errors"

	"github.com/tradeface/suggest_service/internal/conf"
)

type Service struct {
	Mongo   *MongoService
	Elastic *ElasticService
}

func New(cfg *conf.Config) (service *Service, err error) {

	mongoService, err := NewMongoService(cfg)
	if err != nil {
		return service, errors.New("failed to connect to db")
	}

	elasticService, err := NewElasticService(cfg)
	if err != nil {
		return service, errors.New("failed to connect to es")
	}

	return &Service{
		Mongo:   mongoService,
		Elastic: elasticService,
	}, nil
}
