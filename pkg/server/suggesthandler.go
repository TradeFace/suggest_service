package server

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tradeface/suggest_service/pkg/authorization"
	"github.com/tradeface/suggest_service/pkg/model"
	"github.com/tradeface/suggest_service/pkg/store"
)

type suggestHandler struct {
	esQueryCache map[string]esQueryCacheItem
	stores       *store.Stores
	auth         *authorization.AuthChecker
}

type esQueryCacheItem struct {
	expire time.Time
	query  string
}

const QUERY_CACHE_SEC = 300

//TODO: run cache clean every x seconds

func NewSuggestHandler(stores *store.Stores, auth *authorization.AuthChecker) *suggestHandler {
	return &suggestHandler{
		esQueryCache: make(map[string]esQueryCacheItem, 0),
		stores:       stores,
		auth:         auth,
	}
}

func (sh *suggestHandler) getQuery(c echo.Context) (string, error) {

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

func (sh *suggestHandler) getDomainQuery(host string) (string, error) {

	if cacheItem, c := sh.esQueryCache[host]; c {
		if cacheItem.expire.After(time.Now()) {
			return cacheItem.query, nil
		}
		delete(sh.esQueryCache, host)
	}

	dbdomain, err := sh.stores.Domain.GetOneWithHost(host)
	if err != nil {
		return "", err
	}
	filters := sh.getBaseFilters(dbdomain)

	query := fmt.Sprintf(`{
        "post_filter": {
            "bool": {
                "must": {
                    "has_child": {
                        "type": "article",
                        "query": {
                            "bool": {
                                "must": []
                            }
                        }
                    }
                }
            }
        },
        "from": 0,
        "size": %%d,
        "query": {
            "function_score": {
                "query": {
                    "has_child": {
                        "type": "article",
                        "query": {
                            "bool": {
                                "must": {
                                    "match": {
                                        "suggest": {
                                            "query": "%%s",
                                            "fuzziness": 1
                                        }
                                    }
                                },
                                "filter": [%s]
                            }
                        },
                        "min_children": 1,
                        "score_mode": "max"
                    }
                },
                "boost_mode": "multiply",
                "functions": [{
                    "field_value_factor": {
                        "field": "ibScore"
                    }
                }]
            }
        }
    }`, filters)

	sh.esQueryCache[host] = esQueryCacheItem{
		expire: time.Now().Add(time.Second * QUERY_CACHE_SEC),
		query:  query,
	}

	return query, nil
}

func (sh *suggestHandler) getBaseFilters(domain *model.Domain) string {

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

func (sh *suggestHandler) getClassificationFilter(domain *model.Domain) string {

	return fmt.Sprintf(`{
		"term": {
			"classification.path": "%s"
		}
	}`, domain.MainClassification)
}

func (sh *suggestHandler) getSupplierFilter(domain *model.Domain) string {

	catalogsStr, err := json.Marshal(domain.Catalogs)
	if err != nil {
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

func (sh *suggestHandler) getAvailabilityFilter(domain *model.Domain) string {

	states := domain.GetSetting("SEARCH", "disabled_states")

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

func (sh *suggestHandler) getStockFilter(domain *model.Domain) string {

	if false == domain.ModuleIsEnabled("STOCK") {
		return ""
	}

	stockOnly := domain.GetSetting("STOCK", "search_only_in_stock")
	if stockOnly.(bool) == false {
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

func (sh *suggestHandler) getPublicOnlyFilter() string {
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