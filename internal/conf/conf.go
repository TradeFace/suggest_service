package conf

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug  bool   `envconfig:"DEBUG"`
	Addr   string `envconfig:"ADDR" default:":8080"`
	Stage  string `envconfig:"STAGE" default:"dev"`
	Branch string `envconfig:"BRANCH"`
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
		return nil, errors.New(err)
	}
	return cfg, nil
}
