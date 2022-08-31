package config

import (
	"github.com/auth0/go-auth0/management"
	"github.com/raynigon/auth0_exporter/v2/collector"
)

func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func (cfg Auth0ExporterConfig) GetCollectorConfig() collector.CollectorConfig {
	client, err := management.New(*cfg.domain, management.WithClientCredentials(*cfg.clientId, *cfg.clientSecret))
	if err != nil {
		panic(err)
	}
	return collector.CollectorConfig{
		Logger: cfg.GetLogger(),
		Domain: cfg.domain,
		Auth0:  client,
	}
}
