package config

import (
	"context"

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
	ctx := context.Background()
	client, err := management.New(cfg.Domain, management.WithClientCredentials(ctx, cfg.ClientId, cfg.ClientSecret))
	if err != nil {
		panic(err)
	}
	return collector.CollectorConfig{
		Logger: cfg.GetLogger(),
		Domain: &cfg.Domain,
		Auth0:  client,
	}
}
