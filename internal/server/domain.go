package server

import "github.com/labstack/echo/v4"

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
