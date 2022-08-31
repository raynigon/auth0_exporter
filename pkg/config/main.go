package config

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

type Auth0ExporterConfig struct {
	listenAddress *string
	metricsPath   *string
	logLevel      *string
	logFormat     *string
	logOutput     *string
	domain        *string
	clientId      *string
	clientSecret  *string
}

func NewAuth0ExporterConfig() Auth0ExporterConfig {
	config := Auth0ExporterConfig{
		listenAddress: kingpin.Flag("web.listen-address", "Address to listen on for web interface and telemetry.").
			Default(":9776").
			Envar("A0E_WEB_LISTEN_ADDRESS").
			String(),
		metricsPath: kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").
			Default("/metrics").
			Envar("A0E_WEB_TELEMETRY_PATH").
			String(),
		logLevel: kingpin.Flag("log.level", "Sets the loglevel. Valid levels are debug, info, warn, error").
			Default("info").
			Envar("A0E_LOG_LEVEL").
			String(),
		logFormat: kingpin.Flag("log.format", "Sets the log format. Valid formats are json and logfmt").
			Default("logfmt").
			Envar("A0E_LOG_FORMAT").
			String(),
		logOutput: kingpin.Flag("log.output", "Sets the log output. Valid outputs are stdout and stderr").
			Default("stdout").
			Envar("A0E_LOG_OUTPUT").
			String(),
		domain: kingpin.Flag("auth0.domain", "Sets the auth0 domain").
			Envar("A0E_AUTH0_DOMAIN").
			String(),
		clientId: kingpin.Flag("auth0.client-id", "Sets the auth0 clientId").
			Envar("A0E_AUTH0_CLIENT_ID").
			String(),
		clientSecret: kingpin.Flag("auth0.client-secret", "Sets the auth0 clientSecret").
			Envar("A0E_AUTH0_CLIENT_SECRET").
			String(),
	}
	kingpin.Version("0.0.1")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	return config
}

func (cfg Auth0ExporterConfig) GetListeningAccess() string {
	return *cfg.listenAddress
}

func (cfg Auth0ExporterConfig) GetMetricsPath() string {
	return *cfg.metricsPath
}
