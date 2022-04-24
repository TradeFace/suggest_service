package controller

import (
	"errors"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/authorization"
	"github.com/tradeface/suggest_service/pkg/store"
)

type UserController struct {
	StoreProvider *store.Provider
	cfg           *conf.Config
}

func (uc *UserController) Login(c echo.Context) error {
	email := c.QueryParam("email")
	password := c.QueryParam("password")

	res, err := uc.StoreProvider.User.GetWithEmail(email)
	if err != nil || len(res) != 1 {
		sendError(c, errors.New("No login"))
		return nil
	}

	if !authorization.CheckPasswordHash(password, res[0].Password) {
		sendError(c, errors.New("No login"))
		return nil
	}

	token, err := authorization.NewJwtWithUser(res[0], uc.cfg.JWTSalt)
	if err != nil {
		sendError(c, errors.New("No login"))
		return nil
	}
	res[0].Token = token

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
