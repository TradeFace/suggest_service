package main

import (
	"io/ioutil"

	"github.com/tradeface/suggest_service/internal/api"
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/elastic"
	"github.com/tradeface/suggest_service/pkg/health"
	"github.com/tradeface/suggest_service/pkg/mongo"
	"github.com/tradeface/suggest_service/pkg/server"
	"github.com/tradeface/suggest_service/pkg/store"

	"github.com/labstack/echo/v4"
	echolog "github.com/labstack/gommon/log"
	"github.com/rs/zerolog/log"
)

func main() {

	// loads configuration from env and configures logger
	cfg, err := conf.NewDefaultConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	elasticClient, err := elastic.NewClient(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to elastic")
	}

	mongoClient, err := mongo.NewClient(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to mongo")
	}

	stores, err := store.New(mongoClient, elasticClient, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	srv, err := server.NewServer(cfg, stores)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bind api")
	}

	e := echo.New()
	// shut up
	e.Logger.SetOutput(ioutil.Discard)
	e.Logger.SetLevel(echolog.OFF)

	e.GET("/health", echo.WrapHandler(health.Handler()))

	api.RegisterHandlers(e, srv)

	err = e.Start(cfg.Addr)
	log.Fatal().Err(err).Msg("Server failed")

}
