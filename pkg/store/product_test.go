package store

import (
	"reflect"
	"testing"

	"github.com/tradeface/suggest_service/pkg/document"
	"github.com/tradeface/suggest_service/pkg/service"
)

func TestNewProductStore(t *testing.T) {
	type args struct {
		esconn *service.ElasticService
	}
	tests := []struct {
		name string
		args args
		want *ProductStore
	}{
		{
			name: "InstanciateProductStore",
			args: args{
				esconn: &service.ElasticService{},
			},
			want: &ProductStore{
				esconn: &service.ElasticService{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProductStore(tt.args.esconn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProductStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProductStore_Search(t *testing.T) {
	type fields struct {
		esconn *service.ElasticService
	}
	type args struct {
		query string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantResults []*document.Product
		wantErr     bool
	}{
		{
			name:   "esServiceNotInstanciated",
			fields: fields{},
			args: args{
				query: "{anything}",
			},
			wantResults: []*document.Product{},
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProductStore{
				esconn: tt.fields.esconn,
			}
			gotResults, err := p.Search(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProductStore.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResults, tt.wantResults) {
				t.Errorf("ProductStore.Search() = %v, want %v", gotResults, tt.wantResults)
			}
		})
	}
}
