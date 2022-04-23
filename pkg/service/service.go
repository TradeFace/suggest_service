package service

import (
	"errors"
)

type Service struct {
	Mongo   *MongoService
	Elastic *ElasticService
}

type Config struct {
	MongoURI     string
	MongoDB      string
	ElasticURI   string
	ElasticIndex string
}

func New(cfg *Config) (service *Service, err error) {

	services := &Service{}

	mongoService, err := createMongo(cfg)
	if err != nil {
		return service, err
	}
	services.Mongo = mongoService

	elasticService, err := creatElastic(cfg)
	if err != nil {
		return service, err
	}
	services.Elastic = elasticService

	return services, nil
}

func createMongo(cfg *Config) (mongoService *MongoService, err error) {

	if cfg.MongoURI == "" {
		return mongoService, err
	}
	mongoService, err = NewMongoService(cfg)
	if err != nil {
		return mongoService, errors.New("failed to connect to db")
	}
	return mongoService, err
}

func creatElastic(cfg *Config) (elasticService *ElasticService, err error) {

	if cfg.ElasticURI == "" {
		return elasticService, err
	}
	elasticService, err = NewElasticService(cfg)
	if err != nil {
		return elasticService, errors.New("failed to connect to es")
	}
	return elasticService, err
}
