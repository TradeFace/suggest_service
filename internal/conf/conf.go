package conf

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug      bool   `envconfig:"DEBUG"`
	Addr       string `envconfig:"ADDR" default:":8888"`
	Stage      string `envconfig:"STAGE" default:"dev"`
	Branch     string `envconfig:"BRANCH"`
	MongoURI   string `envconfig:"MONGO_URI" default:"mongodb://root:example@localhost:27017/"`
	MongoDB    string `envconfig:"MONGO_DB" default:"suggest_test"`
	ElasticURI string `envconfig:"ELASTIC_URI" default:"http://localhost:9200/"`
	ElasticDB  string `envconfig:"ELASTIC_DB" default:"suggest_testa"`
	// PGDatasource         string `envconfig:"PGDATASOURCE"`
	// OpenIDProvider       string `envconfig:"OPENID_PROVIDER_URL"`
	// ClientID             string `envconfig:"OAUTH_CLIENT_ID"`
	// MetricsWriteInterval int    `envconfig:"METRICS_WRITE_INTERVAL"`
	// DbSecrets            string `envconfig:"DB_SECRET"`
}

func NewDefaultConfig() (*Config, error) {

	cfg := new(Config)
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
