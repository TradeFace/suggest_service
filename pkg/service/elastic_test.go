package service

import (
	"reflect"
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
)

func TestNewElasticService(t *testing.T) {
	type args struct {
		cfg *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *ElasticService
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewElasticService(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewElasticService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewElasticService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestElasticService_startClient(t *testing.T) {
	type fields struct {
		Client    *elasticsearch.Client
		URI       []string
		IndexName string
		User      string
		Password  string
		debug     bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &ElasticService{
				Client:    tt.fields.Client,
				URI:       tt.fields.URI,
				IndexName: tt.fields.IndexName,
				User:      tt.fields.User,
				Password:  tt.fields.Password,
				debug:     tt.fields.debug,
			}
			if err := es.startClient(); (err != nil) != tt.wantErr {
				t.Errorf("ElasticService.startClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestElasticService_Search(t *testing.T) {
	type fields struct {
		Client    *elasticsearch.Client
		URI       []string
		IndexName string
		User      string
		Password  string
		debug     bool
	}
	type args struct {
		query string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantR   *ElasticResult
		wantErr bool
	}{
		{
			name:   "QueryWithoutInstance",
			fields: fields{},
			args: args{
				query: `{"query":{"match_all":{}}}`,
			},

			wantR:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &ElasticService{
				Client:    tt.fields.Client,
				URI:       tt.fields.URI,
				IndexName: tt.fields.IndexName,
				User:      tt.fields.User,
				Password:  tt.fields.Password,
				debug:     tt.fields.debug,
			}
			gotR, err := es.Search(tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("ElasticService.Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("ElasticService.Search() = %v, want %v", gotR, tt.wantR)
			}
		})
	}
}
