package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/tradeface/suggest_service/internal/provider"
)

type DomainController struct {
	StoreProvider *provider.StoreProvider
}

func (dc *DomainController) GetList(c echo.Context) error {

	host := c.QueryParam("filter[host]")
	res, err := dc.StoreProvider.Domain.GetWithHost(host)
	return Output(c, res, err)
}

func (dc *DomainController) Get(c echo.Context) error {

	id := c.Param("id")
	res, err := dc.StoreProvider.Domain.GetWithId(id)
	return Output(c, res, err)
}
