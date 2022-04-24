package controller

import (
	"log"
	"net/http"

	"github.com/google/jsonapi"
	"github.com/labstack/echo/v4"
	"github.com/tradeface/suggest_service/internal/provider"
	"github.com/tradeface/suggest_service/internal/suggest"
)

type ProductController struct {
	StoreProvider *provider.StoreProvider
	queryBuilder  *suggest.QueryBuilder
}

func NewProductController(storeProvider *provider.StoreProvider) *ProductController {
	return &ProductController{
		StoreProvider: storeProvider,
		queryBuilder:  suggest.NewQueryBuilder(storeProvider),
	}
}

func (pc *ProductController) GetList(c echo.Context) error {

	user, err := GetAuthUser(c)
	if err == nil {
		if user.HasRole("ADMIN") {
			log.Println("you are admin")
		}

		Id, err := user.GetClaim("Id")
		log.Println(user, err)
		log.Println(Id)
	}

	query, err := pc.queryBuilder.GetQuery(c)

	if err != nil {
		sendError(c, err)
		return nil
	}

	resProd, err := pc.StoreProvider.Product.Search(query)
	if err != nil {
		sendError(c, err)
		return nil
	}

	payload, err := jsonapi.Marshal(resProd)
	if err != nil {
		sendError(c, err)
		return nil
	}
	return c.JSON(http.StatusOK, payload)
}
