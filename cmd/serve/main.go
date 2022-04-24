package main

import (
	"io/ioutil"

	"github.com/labstack/echo/v4"

	echolog "github.com/labstack/gommon/log"
	"github.com/rs/zerolog/log"
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/internal/provider"
	"github.com/tradeface/suggest_service/internal/server"
	"github.com/tradeface/suggest_service/pkg/middleware"
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

	providers := provider.NewProvider(cfg)

	srv, err := server.NewServer(cfg, providers)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to bind api")
	}

	e := echo.New()
	// shut up
	e.Logger.SetOutput(ioutil.Discard)
	e.Logger.SetLevel(echolog.OFF)

	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	e.Use(middleware.JWTWithConfig(providers.Store, cfg.JWTSalt))

	srv.RegisterHandlers(e)

	log.Info().Str("addr", cfg.Addr).Msg("starting http listener")
	err = e.Start(cfg.Addr)
	log.Fatal().Err(err).Msg("Server failed")
}
