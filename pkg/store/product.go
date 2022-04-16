package store

import (
	"github.com/tradeface/suggest_service/internal/api"
	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/elastic"
)

type ProductStore struct {
	es  *elastic.ElasticClient
	cfg *conf.Config
}

// NewDomain new project store
func NewProductStore(es *elastic.ElasticClient, cfg *conf.Config) ProductStore {
	return ProductStore{es: es, cfg: cfg}
}

func (ds *ProductStore) Search(domain map[string]interface{}, pageNumber int64, pageSize int64, text string, host string) api.ProductPage {

	// str := "dit is een product"
	// attr := Attributes{
	// 	Description: "hahahahahahahhaha",
	// }
	row := map[string]interface{}{
		"id":          "hahaha",
		"description": "asdasdasdasdasd",
	}
	x := "hahaha desc"
	return api.ProductPage{
		Data: []api.Product{
			api.Product{
				// Attributes: {Description: "hahaha desc"},
				Id: row["id"].(api.Id),
				// Links: &struct{Self *github.com/tradeface/suggest_service/internal/api.Link "json:\"self,omitempty\""}{},
				Type: "prod",
			},
			api.Product{
				Attributes: &struct {
					Description *string "json:\"description,omitempty\""
				}{&x},
				Id: "hahaha2",
				// Links: &struct{Self *github.com/tradeface/suggest_service/internal/api.Link "json:\"self,omitempty\""}{},
				Type: "prod",
			},
		},
	}

	// tmpProd.Id = func(i string) *string { return &i }(row["id"])
	// tmpProd.Attributes{
	// 	Description: func(i string) *string { return &i }(row["description"]),
	// }

	// 	Id:   row["id"].(string),
	// 	Type: "product",
	// 	// Attributes: attr,
	// }
	// tmpProd.Id = "hahahaha"
	// tmpProd.Type = "product"
	// tmpProd.Attributes.Description = &str

	// tmp := make([]api.Product, 0)
	// tmp = append(tmp, tmpProd)
	// return tmpProd
}
