package provider

import (
	"github.com/rs/zerolog/log"
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/service"
	"github.com/tradeface/suggest_service/pkg/store"
)

type Provider struct {
	Service *service.Provider
	Store   *store.Provider
}

func NewProvider(cfg *conf.Config) *Provider {

	serviceProvider, err := service.NewProvider(&service.Config{
		MongoURI:     cfg.MongoURI,
		MongoDB:      cfg.MongoDB,
		ElasticURI:   cfg.ElasticURI,
		ElasticIndex: cfg.ElasticIndex,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	storeConf := &store.Config{
		Service: serviceProvider,
	}

	storeProvider, err := store.NewProvider(storeConf)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	return &Provider{
		Service: serviceProvider,
		Store:   storeProvider,
	}
}
