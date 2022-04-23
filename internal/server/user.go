package server

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
)

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
