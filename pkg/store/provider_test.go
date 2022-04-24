package store

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/tradeface/suggest_service/pkg/service"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestNewProvider(t *testing.T) {
	type args struct {
		cfg *Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Provider
		wantErr bool
	}{
		{
			name:    "InstanciateProviderWithoutConfig",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "InstanciateProviderWithoutElasticConfig",
			args: args{
				cfg: &Config{
					Service: &service.Provider{
						Mongo: &service.MongoService{
							Client:   &mongo.Client{},
							Database: &mongo.Database{},
							Ctx:      nil,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "InstanciateProviderWithoutMongoConfig",
			args: args{
				cfg: &Config{
					Service: &service.Provider{
						Elastic: &service.ElasticService{
							Client:    &elasticsearch.Client{},
							URI:       []string{},
							IndexName: "",
							User:      "",
							Password:  "",
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		/*	{
			name: "InstanciateProviderWithConfig",
			args: args{
				cfg: &Config{
					Service: &service.Provider{
						Mongo: &service.MongoService{
							Client:   &mongo.Client{},
							Database: &mongo.Database{},
							Ctx:      nil,
						},
						Elastic: &service.ElasticService{
							Client:    &elasticsearch.Client{},
							URI:       []string{},
							IndexName: "",
							User:      "",
							Password:  "",
						},
					},
				},
			},
			want: &Provider{
				Product: &ProductStore{
					esconn: &service.ElasticService{},
				},
				Domain: &DomainStore{
					dbconn:   &service.MongoService{},
					collName: "domain",
				},
				User: &UserStore{
					dbconn:   &service.MongoService{},
					collName: "user",
				},
				Auth: &AuthStore{
					authUser: make(map[string]*authorization.AuthUser),
				},
				ElasticQuery: &ElasticQueryStore{
					elasticQuery: make(map[string]ElasticQuery),
				},
			},
			wantErr: false,
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewProvider(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProvider() = %v, want %v", got, tt.want)
				fmt.Println(got.Auth)
				fmt.Println(tt.want.Auth)
				fmt.Printf("got %T\n", got.Auth)
				fmt.Printf("want %T\n", tt.want.Auth)
			}
		})
	}
}
