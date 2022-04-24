package server

import (
	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/internal/controller"
	"github.com/tradeface/suggest_service/internal/provider"
)

type Server struct {
	controller *controller.Provider
}

func NewServer(cfg *conf.Config, providers *provider.Provider) (*Server, error) {

	controllerProvider, err := controller.NewProvider(cfg, providers.Store)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to db")
	}
	return &Server{
		controller: controllerProvider,
	}, nil
}

func (srv *Server) RegisterHandlers(e *echo.Echo) {

	//http://localhost:8080/product?filter[host]=www.ib.nl&text=gips
	e.GET("/product", srv.controller.Product.GetList)

	//http://localhost:8080/domain/?filter[host]=www.ib.nl
	e.GET("/domain", srv.controller.Domain.GetList)

	//http://localhost:8080/domain/537e3ea78812e9f0e7331733
	e.GET("/domain/:id", srv.controller.Domain.Get)
}
