package main

import (
	"io/ioutil"

	"tradeface.nl/internal/conf"
	"tradeface.nl/pkg/health"

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

	e := echo.New()
	// shut up
	e.Logger.SetOutput(ioutil.Discard)
	e.Logger.SetLevel(echolog.OFF)

	e.GET("/health", echo.WrapHandler(health.Handler()))

	err = e.Start(cfg.Addr)
	log.Fatal().Err(err).Msg("Server failed")

}
