package conf

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// ErrMissingEnvironmentStage missing stage configuration
	ErrMissingEnvironmentStage = errors.New("Missing Stage ENV Variable")

	// ErrMissingEnvironmentBranch missing branch configuration
	ErrMissingEnvironmentBranch = errors.New("Missing Branch ENV Variable")
)

// Config for the environment
type Config struct {
	Debug        bool   `envconfig:"DEBUG"`
	Addr         string `envconfig:"ADDR" default:":8888"`
	Stage        string `envconfig:"STAGE" default:"dev"`
	Branch       string `envconfig:"BRANCH"`
	ElasticURI   string `envconfig:"ELASTIC_URI" default:"http://127.0.0.1:9200"`
	ElasticIndex string `envconfig:"ELASTIC_INDEX" default:"suggest_testa"`
	MongoURI     string `envconfig:"MONGO_URI" default:"mongodb://root:example@127.0.0.1:27017"`
	MongoDB      string `envconfig:"MONGO_DB" default:"suggest_test"`
	JWTSalt      string `envconfig:"JWT_SALT" default:"abc12345"`
}

func (cfg *Config) validate() error {
	if cfg.Stage == "" {
		return ErrMissingEnvironmentStage
	}
	// if cfg.Branch == "" {
	// 	return ErrMissingEnvironmentBranch
	// }

	return nil
}

func (cfg *Config) logging() error {

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	if cfg.Stage == "local" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return nil
}

// NewDefaultConfig reads configuration from environment variables and validates it
func NewDefaultConfig() (*Config, error) {
	cfg := new(Config)
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to parse environment config")
	}

	err = cfg.validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed validation of config")
	}
	err = cfg.logging()
	if err != nil {
		return nil, errors.Wrap(err, "failed setup logging based on config")
	}
	log.Info().Str("stage", cfg.Stage).Bool("debug", cfg.Debug).Msg("logging configured")
	log.Info().Str("stage", cfg.Stage).Str("branch", cfg.Branch).Msg("Configuration loaded")

	return cfg, nil
}
