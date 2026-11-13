package conf

import (
	"net"
	"net/url"
	"strconv"

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
	Auth     AuthConf
	Tracing  tracing.Config
}

type AuthConf struct {
	JWKSURL        string
	Issuer         string
	Audience       string
	SSOProviderID  string
	AllowedDomain  string
	AllowedOrigins string
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
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(r.Username, r.Password),
		Host:     net.JoinHostPort(r.Host, strconv.Itoa(r.Port)),
		Path:     "/" + r.DB,
		RawQuery: url.Values{"sslmode": {r.SSLMode}}.Encode(),
	}

	return u.String()
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
