package main

import (
	"io/ioutil"

	"github.com/labstack/echo/v4"

	echolog "github.com/labstack/gommon/log"
	"github.com/rs/zerolog/log"
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/internal/server"
	"github.com/tradeface/suggest_service/pkg/middleware"
	"github.com/tradeface/suggest_service/pkg/service"
	"github.com/tradeface/suggest_service/pkg/store"
)

//TODO: config cli/dockersecrets

const (
	// APPNAME contains the name of the program
	APPNAME = "suggest_service"
	// APPVERSION contains the version of the program
	APPVERSION = "0.0.3"
)

func main() {

	cfg, err := conf.NewDefaultConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	srvConf := &service.Config{
		MongoURI:     cfg.MongoURI,
		MongoDB:      cfg.MongoDB,
		ElasticURI:   cfg.ElasticURI,
		ElasticIndex: cfg.ElasticIndex,
	}

	services, err := service.New(srvConf)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}

	stores, err := store.New(services)
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

	// tmp Login route
	// e.GET("/login", login)

	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	e.Use(middleware.JWTWithConfig(&middleware.JWTConfig{}, stores.Auth))

	srv.RegisterHandlers(e)

	log.Info().Str("addr", cfg.Addr).Msg("starting http listener")
	err = e.Start(cfg.Addr)
	log.Fatal().Err(err).Msg("Server failed")
}
