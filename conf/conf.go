package conf

import (
	"fmt"
	"net/url"

	"github.com/wahed-tech/go-lib/v2/cfg"
	"github.com/wahed-tech/go-lib/v2/log"
	"github.com/wahed-tech/go-lib/v2/tracing"
)

var R RootConf

type RootConf struct {
	App      cfg.App
	Logger   log.Logger
	Server   cfg.Server
	Kafka    KafkaConf
	Postgres PostgresConf
	Tracing  tracing.Config
}

type PostgresConf struct {
	Username      string
	Password      string
	Host          string
	Port          int
	DB            string
	SSLMode       string
	MigrationsDir string
}

func (r PostgresConf) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		url.QueryEscape(r.Username), url.QueryEscape(r.Password), r.Host, r.Port, r.DB, r.SSLMode)
}

type KafkaConf struct {
	ConfigMap           map[string]any
	Topics              KafkaTopicsConf
	ProcessingTimeoutMs int
	BootstrapServers    string
}

type KafkaTopicsConf struct {
	ExampleEvents string
}

func Init() {
	cfg.Setup("banking-ops-svc", &R)
	R.Kafka.ConfigMap["bootstrap.servers"] = R.Kafka.BootstrapServers
}
