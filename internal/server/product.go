package server

import (
	"log"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
)

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
