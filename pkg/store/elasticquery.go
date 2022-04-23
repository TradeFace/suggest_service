package store

import (
	"errors"
	"time"
)

type ElasticQueryStore struct {
	elasticQuery map[string]ElasticQuery
}

type ElasticQuery struct {
	query  string
	expire time.Time
}

func NewElasticQueryStore() *ElasticQueryStore {

	elasticQueryStore := &ElasticQueryStore{}

	done := make(chan bool)
	elasticQueryStore.forever(done)

	return elasticQueryStore
}

func (eq *ElasticQueryStore) AddQuery(key string, query string, cacheMinutes int) {

	eq.elasticQuery[key] = ElasticQuery{
		query:  query,
		expire: time.Now().Add(time.Duration(cacheMinutes) * time.Minute),
	}
}

func (eq *ElasticQueryStore) GetQuery(key string) (query string, err error) {

	if esQuery, ok := eq.elasticQuery[key]; ok {
		if esQuery.expire.Before(time.Now()) {
			delete(eq.elasticQuery, key)
			return "", errors.New("query not found")
		}
		return esQuery.query, nil
	}
	return "", errors.New("query not found")
}

func (eq *ElasticQueryStore) forever(done <-chan bool) {

	//remove all expired elastic queries
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-done:
				ticker.Stop()
				return
			case <-ticker.C:
				for key, val := range eq.elasticQuery {
					if val.expire.Before(time.Now()) {
						delete(eq.elasticQuery, key)
					}
				}
			}
		}
	}()
}
