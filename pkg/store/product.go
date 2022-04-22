package store

import (
	"fmt"
	"log"

	"github.com/tradeface/suggest_service/internal/conf"
	"github.com/tradeface/suggest_service/pkg/elastic"
	"github.com/tradeface/suggest_service/pkg/model"
)

type ProductStore struct {
	es  *elastic.ElasticClient
	cfg *conf.Config
}

// NewDomain new project store
func NewProductStore(es *elastic.ElasticClient, cfg *conf.Config) ProductStore {

	return ProductStore{es: es, cfg: cfg}
}

func (ds *ProductStore) Search(domain map[string]interface{}, pageNumber int64, pageSize int64, text string, host string) (products []*model.Product) {

	// ds.es.Post("blah")

	res := ds.es.Search(text, model.Product{})
	log.Println("-----search res------")
	log.Println(res)
	log.Println("-----search res end------")
	for _, hit := range res.Hits.Hits {
		log.Println("-----search res hit------")
		log.Println(hit.Source)
		log.Printf("%T", hit.Source)
		// p := hit.Source.(model.Product)
		// temporaryVariable, _ := json.Marshal(hit.Source)
		// p := model.Product{}
		// err := json.Marshal(temporaryVariable, &p)
		// if err != nil {
		// 	// Catch the exception to handle it as per your need
		// }
		// switch hit.Source.(type) {
		// case *interface{}:
		// 	p := hit.Source.(*interface{}).(*model.Product)
		// 	log.Println(p)
		// }

		products = append(products, hit.Source)
		log.Println("-----search res hit end------")
		// 	// p := &model.Product{}
	}
	// for _, hit := range res["hits"].(map[string]interface{})["hits"].([]interface{}) {
	// 	log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	// 	p := &model.Product{
	// 		ID:          hit.(map[string]interface{})["_id"].(string),
	// 		Description: hit.(map[string]interface{})["_source"].(map[string]interface{})["description"].(string),
	// 	}
	// 	products = append(products, p)
	// }
	// str := "dit is een product"
	// attr := Attributes{
	// 	Description: "hahahahahahahhaha",
	// }
	// row := map[string]interface{}{
	// 	"id":          "hahaha",
	// 	"description": "asdasdasdasdasd",
	// }
	// x := "hahaha desc"
	// return []model.Product{
	// 		api.Product{
	// 			// Attributes: {Description: "hahaha desc"},
	// 			Id: "hahaha1",
	// 			// Links: &struct{Self *github.com/tradeface/suggest_service/internal/api.Link "json:\"self,omitempty\""}{},
	// 			Type: "prod",
	// 		},
	// 		api.Product{
	// 			Attributes: &struct {
	// 				Description *string "json:\"description,omitempty\""
	// 			}{&x},
	// 			Id: "hahaha2",
	// 			// Links: &struct{Self *github.com/tradeface/suggest_service/internal/api.Link "json:\"self,omitempty\""}{},
	// 			Type: "prod",
	// 		},
	// 	},
	// }
	// products = append(products, ds.fixtureProduct(1))
	// products = append(products, ds.fixtureProduct(2))
	return products

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

func (ds *ProductStore) fixtureProduct(id int) *model.Product {
	return &model.Product{
		ID:          fmt.Sprintf("%d", id),
		Description: "hahahhahh 1",
	}
}
