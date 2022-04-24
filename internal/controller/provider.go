package controller

import (
	"errors"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"github.com/tradeface/jwt_service/pkg/authorization"
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/store"
)

type Provider struct {
	Domain  *DomainController
	Product *ProductController
	User    *UserController
}

func NewProvider(cfg *conf.Config, storeProvider *store.Provider) (*Provider, error) {

	productController := NewProductController(storeProvider)
	return &Provider{
		Domain: &DomainController{
			StoreProvider: storeProvider,
		},
		User: &UserController{
			StoreProvider: storeProvider,
			cfg:           cfg,
		},
		Product: productController,
	}, nil
}

func GetAuthUser(c echo.Context) (*authorization.AuthUser, error) {

	user := c.Get("authUser")
	if user == nil {
		return nil, errors.New("no auth user available")
	}
	return user.(*authorization.AuthUser), nil
}

func Output(c echo.Context, res interface{}, err error) error {

	if err != nil {
		sendError(c, err)
		return nil
	}
	payload, err := jsonapi.Marshal(res)
	if err != nil {
		sendError(c, err)
		return nil
	}
	return c.JSON(http.StatusOK, payload)
}

func sendError(c echo.Context, err error) {

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusBadRequest)
	jsonapi.MarshalErrors(c.Response().Writer, []*jsonapi.ErrorObject{{
		Title:  "Request Error",
		Detail: err.Error(),
		Status: "400",
	}})
}
