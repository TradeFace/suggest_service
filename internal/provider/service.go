package provider

// Not in use, demo purpose

// import (
// 	"errors"

// 	"github.com/tradeface/suggest_service/internal/conf"
// 	"github.com/tradeface/suggest_service/pkg/service"
// )

// type ServiceProvider struct {
// 	Mongo   *service.MongoService
// 	Elastic *service.ElasticService
// }

// func NewServiceProvider(cfg *conf.Config) (serviceProvider *ServiceProvider, err error) {

// 	serviceProvider = &ServiceProvider{}

// 	mongoService, err := service.NewMongoService(&service.Config{
// 		MongoURI: cfg.MongoURI,
// 		MongoDB:  cfg.MongoDB,
// 	})
// 	if err != nil {
// 		return serviceProvider, errors.New("failed to connect to db")
// 	}
// 	serviceProvider.Mongo = mongoService

// 	elasticService, err := service.NewElasticService(&service.Config{
// 		ElasticURI:   cfg.ElasticURI,
// 		ElasticIndex: cfg.ElasticIndex,
// 	})
// 	if err != nil {
// 		return serviceProvider, errors.New("failed to connect to es")
// 	}
// 	serviceProvider.Elastic = elasticService

// 	return serviceProvider, nil
// }
