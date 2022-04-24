package suggest

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/tradeface/suggest_service/pkg/document"
	"github.com/tradeface/suggest_service/pkg/store"
)

const DOMAIN_QUERY_CACHE_MIN = 5

type QueryBuilder struct {
	stores *store.Provider
}

func NewQueryBuilder(storeProvider *store.Provider) *QueryBuilder {
	return &QueryBuilder{
		stores: storeProvider,
	}
}

func (sh *QueryBuilder) GetQuery(c echo.Context) (string, error) {

	text := c.QueryParam("text")
	host := c.QueryParam("filter[host]")
	pageSize, err := strconv.Atoi(c.QueryParam("page[size]"))
	if err != nil {
		pageSize = 3
	}
	if pageSize > 10 || pageSize < 1 {
		return "", fmt.Errorf("page[size] out of range: 1 - 10")
	}

	q, err := sh.getDomainQuery(host)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(q, pageSize, text), nil
}

func (sh *QueryBuilder) getDomainQuery(host string) (string, error) {

	query, err := sh.stores.ElasticQuery.GetQuery("domain_suggest_" + host)
	if err == nil {
		return query, nil
	}

	// dbdomain, err := sh.stores.Domain.GetOneWithHost(host)
	// if err != nil {
	// 	return "", err
	// }
	// filters := sh.getBaseFilters(dbdomain)

	query = fmt.Sprintf(`{
        
        "from": 0,
        "size": %%d,
    
		"query": {
			"bool": {
				"must": {
					"match": {
						"description": {
							"query": "%%s",
							"fuzziness": 1
						}
					}
				}
				
			}
		}
                       
                        
    }`) //"filter": [%s] , filters

	sh.stores.ElasticQuery.AddQuery("domain_suggest_"+host, query, DOMAIN_QUERY_CACHE_MIN)

	return query, nil
}

func (sh *QueryBuilder) getBaseFilters(domain *document.Domain) string {

	filters := make([]string, 0)
	if res := sh.getClassificationFilter(domain); res != "" {
		filters = append(filters, res)
	}

	if res := sh.getSupplierFilter(domain); res != "" {
		filters = append(filters, res)
	}

	if res := sh.getAvailabilityFilter(domain); res != "" {
		filters = append(filters, res)
	}

	if res := sh.getStockFilter(domain); res != "" {
		filters = append(filters, res)
	}

	if res := sh.getPublicOnlyFilter(); res != "" {
		filters = append(filters, res)
	}

	return strings.Join(filters[:], ",")
}

func (sh *QueryBuilder) getClassificationFilter(domain *document.Domain) string {

	return fmt.Sprintf(`{
		"term": {
			"classification.path": "%s"
		}
	}`, domain.MainClassification)
}

func (sh *QueryBuilder) getSupplierFilter(domain *document.Domain) string {

	catalogsStr, err := json.Marshal(domain.Catalogs)
	if err != nil || len(domain.Catalogs) == 0 {
		return ""
	}

	return fmt.Sprintf(`{
		"bool": {
			"should": {
				"terms": {
					"supplier": %s
				}
			}
		}
	}`, catalogsStr)
}

func (sh *QueryBuilder) getAvailabilityFilter(domain *document.Domain) string {

	states, err := domain.GetSetting("SEARCH", "disabled_states")
	if err != nil {
		return ""
	}

	statesStr := make([]string, 0)
	for _, val := range states.(map[string]interface{}) {
		statesStr = append(statesStr, val.(string))
	}

	state, err := json.Marshal(statesStr)
	if err != nil {
		return ""
	}

	return fmt.Sprintf(`{
		"bool": {
			"must_not": {
				"terms": {
					"availability": %s
				}
			}
		}
	}`, state)
}

func (sh *QueryBuilder) getStockFilter(domain *document.Domain) string {

	if !domain.ModuleIsEnabled("STOCK") {
		return ""
	}

	stockOnly, err := domain.GetSetting("STOCK", "search_only_in_stock")
	if err != nil {
		return ""
	}

	if !stockOnly.(bool) {
		return ""
	}
	return `{
		"bool": {
			"should": {
				"terms": {
					"inStock": 1
				}
			}
		}
	}`
}

func (sh *QueryBuilder) getPublicOnlyFilter() string {
	// Only allow public articles to return or articles that this user may see.
	// E.g. Filter out private articles for non IB admins.
	// if ($this->authChecker->isGranted('ROLE_MAY_SEE_NON_PUBLIC_ARTICLES')) {
	//     return;
	// }

	// $companyId = $this->userContainer->getCompanyIdOrNull();
	//TODO: get grant and companyId when JWT available
	companyId := ""

	if companyId == "" {

		return `{
			"term": {
				"isPublic": true
			}
		}`
	}

	return fmt.Sprintf(`{
		"bool": {
			"should": [
				{
					"term": {
						"isPublic": true
					}
				},
				{
					"term": {
						"visibleFor": "%s"
					}
				}
			]
		}
	}`, companyId)

}
