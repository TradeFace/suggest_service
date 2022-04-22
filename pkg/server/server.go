package server

import (
	"net/http"

	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/authorization"
	"github.com/tradeface/suggest_service/pkg/store"
)

type Server struct {
	cfg            *conf.Config
	stores         *store.Stores
	suggestHandler *suggestHandler
	auth           *authorization.AuthChecker
}

func NewServer(cfg *conf.Config, stores *store.Stores, auth *authorization.AuthChecker) (*Server, error) {

	suggestHandler := NewSuggestHandler(stores, auth)
	return &Server{
		cfg:            cfg,
		stores:         stores,
		suggestHandler: suggestHandler,
		auth:           auth,
	}, nil
}

func (srv *Server) RegisterHandlers(e *echo.Echo) {

	//http://localhost:8080/product?filter[host]=www.ib.nl&text=gips
	e.GET("/product", srv.GetProductList)

	//http://localhost:8080/domain/?filter[host]=www.ib.nl
	e.GET("/domain", srv.GetDomainList)

	//http://localhost:8080/domain/537e3ea78812e9f0e7331733
	e.GET("/domain/:id", srv.GetDomain)
}

func (srv *Server) GetProductList(c echo.Context) error {

	query, err := srv.suggestHandler.getQuery(c)

	if err != nil {
		srv.sendError(c, err)
		return nil
	}

	resProd, err := srv.stores.Product.Search(query)
	if err != nil {
		srv.sendError(c, err)
		return nil
	}

	payload, err := jsonapi.Marshal(resProd)
	if err != nil {
		srv.sendError(c, err)
		return nil
	}
	return c.JSON(http.StatusOK, payload)
}

func (srv *Server) GetDomainList(c echo.Context) error {

	host := c.QueryParam("filter[host]")
	res, err := srv.stores.Domain.GetWithHost(host)
	if err != nil {
		srv.sendError(c, err)
		return nil
	}

	payload, err := jsonapi.Marshal(res)
	if err != nil {
		srv.sendError(c, err)
		return nil
	}
	return c.JSON(http.StatusOK, payload)
}

func (srv *Server) GetDomain(c echo.Context) error {

	id := c.Param("id")
	res, err := srv.stores.Domain.GetWithId(id)
	if err != nil {
		srv.sendError(c, err)
		return nil
	}

	payload, err := jsonapi.Marshal(res)
	if err != nil {
		srv.sendError(c, err)
		return nil
	}
	return c.JSON(http.StatusOK, payload)
}

func (srv *Server) sendError(c echo.Context, err error) {

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusBadRequest)
	jsonapi.MarshalErrors(c.Response().Writer, []*jsonapi.ErrorObject{{
		Title:  "Request Error",
		Detail: err.Error(),
		Status: "400",
	}})
}
