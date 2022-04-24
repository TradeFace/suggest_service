package controller

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/tradeface/suggest_service/internal/provider"
)

type UserController struct {
	StoreProvider *provider.StoreProvider
}

func (uc *UserController) Login(c echo.Context) error {
	email := c.QueryParam("email")
	password := c.QueryParam("password")
	res, err := uc.StoreProvider.User.Login(email, password)
	return Output(c, res, err)
}

func (uc *UserController) Get(c echo.Context) error {

	user, err := GetAuthUser(c)
	if err != nil {
		err := fmt.Errorf("not allowed")
		return Output(c, "", err)
	}
	if user.HasRole("ADMIN") {
		log.Println("you are admin")
	}
	claimId, err := user.GetClaim("Id")
	if err != nil {
		sendError(c, err)
		return nil
	}

	id := c.Param("id")
	if id != claimId.(string) {
		err := fmt.Errorf("not allowed")
		return Output(c, "", err)
	}
	res, err := uc.StoreProvider.User.GetWithId(id)
	return Output(c, res, err)
}
