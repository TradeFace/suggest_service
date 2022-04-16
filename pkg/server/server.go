package server

import (
	"log"

	"github.com/tradeface/suggest_service/internal/api"
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/store"

	"net/http"

	"github.com/labstack/echo/v4"
)

// Server represents all server handlers.
type Server struct {
	cfg    *conf.Config
	stores *store.Stores
}

// NewServer new api server
func NewServer(cfg *conf.Config, stores *store.Stores) (*Server, error) {
	return &Server{cfg: cfg, stores: stores}, nil
}

func (sv *Server) GetProducts(ctx echo.Context, params api.GetProductsParams) error {

	pageNumber, pageSize, text, host := sv.getProductsParams(params)

	log.Println(pageNumber)
	log.Println(pageSize)
	log.Println(text)
	log.Println(host)
	domainData := sv.stores.Domain.GetByHost(host)
	log.Println(domainData)
	productData := sv.stores.Product.Search(domainData, pageNumber, pageSize, text, host)

	// return ctx.JSON(http.StatusOK, domainData)
	return ctx.JSON(http.StatusOK, &productData)
}

func (sv *Server) getProductsParams(params api.GetProductsParams) (int64, int64, string, string) {

	var pageNumber int64 = 0
	var pageSize int64 = 0
	var text string

	if params.PageNumber != nil {
		pageNumber = int64(*params.PageNumber)
	}
	if params.PageSize != nil {
		pageSize = int64(*params.PageSize)
	}
	if params.Text != nil {
		text = string(*params.Text)
	}

	return pageNumber, pageSize, text, string(params.FilterHost)
}

// GetProject (GET /projects/{id})
func (sv *Server) GetProduct(ctx echo.Context, id string) error {

	// resProj, err := sv.stores.Products.GetByID(ctx.Request().Context())
	// if err != nil {
	// 	if _, ok := err.(*store.ProjectNotFoundError); ok {
	// 		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	// 	}
	// 	return err
	// }

	return ctx.JSON(http.StatusOK, "not implemented")
}
