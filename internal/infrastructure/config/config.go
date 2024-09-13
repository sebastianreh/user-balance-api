package config

import (
	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		ProjectName    string `default:"user-balance-api"`
		ProjectVersion string `envconfig:"PROJECT_VERSION" default:"0.0.1"`
		Port           string `envconfig:"PORT" default:"8000" required:"true"`
		Prefix         string `envconfig:"PREFIX" default:"/user-balance-api"`
		Env            string `envconfig:"ENV" default:"prod"`
		Postgres       struct {
			Host     string `envconfig:"POSTGRES_HOST" default:"127.0.0.1"`
			Port     string `envconfig:"POSTGRES_PORT" default:"5432"`
			User     string `envconfig:"POSTGRES_USER" default:"postgres"`
			Password string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
			DbName   string `envconfig:"POSTGRES_DATABASE" default:"user-balance"`
		}
	}
)

var (
	Configs Config
)

func NewConfig() Config {
	if err := envconfig.Process("", &Configs); err != nil {
		panic(err.Error())
	}

	return Configs
}
