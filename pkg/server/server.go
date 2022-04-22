package server

import (
	"errors"
	"fmt"
	"log"
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

	e.GET("/user/login", srv.LoginUser)

	//http://localhost:8888/user/6262ce0dafd1acb9dfbc4f87
	e.GET("/user/:id", srv.GetUser)
}

func (srv *Server) GetProductList(c echo.Context) error {

	user, err := srv.GetAuthUser(c)
	if err == nil {
		if user.HasRole("ADMIN") {
			log.Println("you are admin")
		}

		Id, err := user.GetClaim("Id")
		log.Println(user, err)
		log.Println(Id)
	}

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
	return srv.Output(c, res, err)
}

func (srv *Server) GetDomain(c echo.Context) error {

	id := c.Param("id")
	res, err := srv.stores.Domain.GetWithId(id)
	return srv.Output(c, res, err)
}

func (srv *Server) LoginUser(c echo.Context) error {
	email := c.QueryParam("email")
	password := c.QueryParam("password")
	res, err := srv.stores.User.Login(email, password)
	return srv.Output(c, res, err)
}

func (srv *Server) GetUser(c echo.Context) error {

	user, err := srv.GetAuthUser(c)
	if err != nil {
		err := fmt.Errorf("not allowed")
		return srv.Output(c, "", err)
	}
	if user.HasRole("ADMIN") {
		log.Println("you are admin")
	}
	claimId, err := user.GetClaim("Id")
	if err != nil {
		srv.sendError(c, err)
		return nil
	}

	id := c.Param("id")
	if id != claimId.(string) {
		err := fmt.Errorf("not allowed")
		return srv.Output(c, "", err)
	}
	res, err := srv.stores.User.GetWithId(id)
	return srv.Output(c, res, err)
}

func (srv *Server) GetAuthUser(c echo.Context) (*authorization.AuthUser, error) {

	user := c.Get("authUser")
	if user == nil {
		return nil, errors.New("no auth user available")
	}
	return user.(*authorization.AuthUser), nil
}

func (srv *Server) Output(c echo.Context, res interface{}, err error) error {

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
