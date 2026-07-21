package config

import (
	"github.com/alecthomas/kong"
)

type Auth0ExporterConfig struct {
	ListenAddress string `name:"web.listen-address" default:":9776" env:"A0E_WEB_LISTEN_ADDRESS" help:"Address to listen on for web interface and telemetry."`
	MetricsPath   string `name:"web.telemetry-path" default:"/metrics" env:"A0E_WEB_TELEMETRY_PATH" help:"Path under which to expose metrics."`
	LogLevel      string `name:"log.level" default:"info" env:"A0E_LOG_LEVEL" help:"Sets the loglevel. Valid levels are debug, info, warn, error"`
	LogFormat     string `name:"log.format" default:"logfmt" env:"A0E_LOG_FORMAT" help:"Sets the log format. Valid formats are json and logfmt"`
	LogOutput     string `name:"log.output" default:"stdout" env:"A0E_LOG_OUTPUT" help:"Sets the log output. Valid outputs are stdout and stderr"`
	Domain        string `name:"auth0.domain" env:"A0E_AUTH0_DOMAIN" help:"Sets the auth0 domain"`
	ClientId      string `name:"auth0.client-id" env:"A0E_AUTH0_CLIENT_ID" help:"Sets the auth0 clientId"`
	ClientSecret  string `name:"auth0.client-secret" env:"A0E_AUTH0_CLIENT_SECRET" help:"Sets the auth0 clientSecret"`
}

func NewAuth0ExporterConfig() Auth0ExporterConfig {
	config := Auth0ExporterConfig{}
	kong.Parse(&config,
		kong.Name("auth0_exporter"),
		kong.UsageOnError(),
	)
	return config
}

func (cfg Auth0ExporterConfig) GetListeningAccess() string {
	return cfg.ListenAddress
}

func (cfg Auth0ExporterConfig) GetMetricsPath() string {
	return cfg.MetricsPath
}
