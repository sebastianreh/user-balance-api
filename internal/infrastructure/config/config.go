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
			Host          string `envconfig:"POSTGRES_HOST" default:"127.0.0.1"`
			Port          string `envconfig:"POSTGRES_PORT" default:"5432"`
			User          string `envconfig:"POSTGRES_USER" default:"postgres"`
			Password      string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
			DBName        string `envconfig:"POSTGRES_DATABASE" default:"user-balance"`
			ReconnectIdle int    `envconfig:"POSTGRES_RECONNECT_IDLE" default:"3"`
		}
		SMTP struct {
			Username string `envconfig:"SMTP_USER" default:"apikey"`
			// Check the author sent email for this value
			Password string `envconfig:"SMTP_PASSWORD" default:""`
			From     string `envconfig:"SMTP_EMAIL" default:"sebastianreh@gmail.com"`
			Host     string `envconfig:"SMTP_HOST" default:"smtp.sendgrid.net"`
			Port     string `envconfig:"SMTP_PORT" default:"587"`
			SendTo   string `envconfig:"SMTP_SEND_TO" default:"sebastianreh@gmail.com"`
		}
		Workers struct {
			MigrationWorkersSize     int `envconfig:"MIGRATION_WORKERS_SIZE" default:"5"`
			MigrationWorkerBatchSize int `envconfig:"MIGRATION_WORKERS_BATCH_SIZE" default:"400"`
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
