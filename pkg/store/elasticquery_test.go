package store

import (
	"reflect"
	"testing"
	"time"
)

func TestNewElasticQueryStore(t *testing.T) {
	tests := []struct {
		name string
		want *ElasticQueryStore
	}{
		{
			name: "InstanciateQueryStore",
			want: &ElasticQueryStore{
				elasticQuery: make(map[string]ElasticQuery),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewElasticQueryStore(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewElasticQueryStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElasticQueryStore_AddQuery(t *testing.T) {
	type fields struct {
		elasticQuery map[string]ElasticQuery
	}
	type args struct {
		key          string
		query        string
		cacheMinutes int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "AddAQuery",
			fields: fields{
				elasticQuery: map[string]ElasticQuery{},
			},
			args: args{
				key:          "testQuery",
				query:        "testQuerytestQuerytestQuerytestQuery",
				cacheMinutes: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eq := &ElasticQueryStore{
				elasticQuery: tt.fields.elasticQuery,
			}
			eq.AddQuery(tt.args.key, tt.args.query, tt.args.cacheMinutes)
		})
	}
}

func TestElasticQueryStore_GetQuery(t *testing.T) {
	type fields struct {
		elasticQuery map[string]ElasticQuery
	}
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantQuery string
		wantErr   bool
	}{
		{
			name:   "GetNoneExistingQuery",
			fields: fields{},
			args: args{
				key: "noneExistingQuery",
			},
			wantQuery: "",
			wantErr:   true,
		},
		{
			name: "GetExpiredQuery",
			fields: fields{
				elasticQuery: map[string]ElasticQuery{
					"expiredQuery": {
						query:  "expiredQueryexpiredQuery",
						expire: time.Now().Add(-time.Second * 10),
					},
				},
			},
			args: args{
				key: "expiredQuery",
			},
			wantQuery: "",
			wantErr:   true,
		},
		{
			name: "GetActiveQuery",
			fields: fields{
				elasticQuery: map[string]ElasticQuery{
					"activeQuery": {
						query:  "activeQueryactiveQuery",
						expire: time.Now().Add(time.Second * 10),
					},
				},
			},
			args: args{
				key: "activeQuery",
			},
			wantQuery: "activeQueryactiveQuery",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eq := &ElasticQueryStore{
				elasticQuery: tt.fields.elasticQuery,
			}
			gotQuery, err := eq.GetQuery(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("ElasticQueryStore.GetQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotQuery != tt.wantQuery {
				t.Errorf("ElasticQueryStore.GetQuery() = %v, want %v", gotQuery, tt.wantQuery)
			}
		})
	}
}

func TestElasticQueryStore_forever(t *testing.T) {
	type fields struct {
		elasticQuery map[string]ElasticQuery
	}
	type args struct {
		done <-chan bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "RunForever",
			fields: fields{
				elasticQuery: map[string]ElasticQuery{},
			},
			args: args{
				done: make(<-chan bool),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eq := &ElasticQueryStore{
				elasticQuery: tt.fields.elasticQuery,
			}
			eq.forever(tt.args.done)
		})
	}
}
