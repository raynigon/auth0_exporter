package collector

import (
	"context"
	"errors"
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/auth0/go-auth0/management"
)

var (
	usersSubsystem = "users"
)

type Auth0UsersCollector struct {
	config        CollectorConfig
	totalUsers    *prometheus.Desc
	blockedUsers  *prometheus.Desc
	emailVerified *prometheus.Desc
}

func init() {
	registerCollector(usersSubsystem, NewAuth0UsersCollector)
}

// NewOrgActionsCollector returns a new Collector exposing actions billing stats.
func NewAuth0UsersCollector(config CollectorConfig, ctx context.Context) (Collector, error) {
	tenant, err := config.Auth0.Tenant.Read()
	if err != nil {
		return nil, err
	}
	if config.Domain == nil {
		return nil, errors.New("domain in CollectorConfig is nil")
	}
	if tenant.FriendlyName == nil {
		return nil, errors.New("FriendlyName in tenant is nil")
	}

	collector := &Auth0UsersCollector{
		config: config,
		totalUsers: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, usersSubsystem, "total"),
			"The number of existing users.",
			nil, prometheus.Labels{"domain": *config.Domain, "tenant": *tenant.FriendlyName},
		),
		blockedUsers: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, usersSubsystem, "blocked"),
			"The number of blocked users.",
			nil, prometheus.Labels{"domain": *config.Domain, "tenant": *tenant.FriendlyName},
		),
		emailVerified: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, usersSubsystem, "email_verified"),
			"The number of users which have a verified email address.",
			nil, prometheus.Labels{"domain": *config.Domain, "tenant": *tenant.FriendlyName},
		),
	}
	err = collector.Reload(ctx)
	if err != nil {
		return nil, err
	}
	return collector, nil
}

func (auc *Auth0UsersCollector) queryUserCount(opts ...management.RequestOption) (int, error) {
	response, err := auc.config.Auth0.User.List(opts...)
	if err != nil {
		return 0, err
	}
	return response.Total, nil
}

func (auc *Auth0UsersCollector) updateTotalUsers(wg *sync.WaitGroup, errors chan error, ch chan<- prometheus.Metric) {
	totalUsers, err := auc.queryUserCount()
	if err != nil {
		errors <- err
		wg.Done()
		return
	}
	ch <- prometheus.MustNewConstMetric(auc.totalUsers, prometheus.GaugeValue, float64(totalUsers))
	wg.Done()
}

func (auc *Auth0UsersCollector) updateBlockedUsers(wg *sync.WaitGroup, errors chan error, ch chan<- prometheus.Metric) {
	blockedUsers, err := auc.queryUserCount(management.Parameter("q", "blocked:true"))
	if err != nil {
		errors <- err
		wg.Done()
		return
	}
	ch <- prometheus.MustNewConstMetric(auc.blockedUsers, prometheus.GaugeValue, float64(blockedUsers))
	wg.Done()
}

func (auc *Auth0UsersCollector) updateEmailVerified(wg *sync.WaitGroup, errors chan error, ch chan<- prometheus.Metric) {
	emailVerified, err := auc.queryUserCount(management.Parameter("q", "email_verified:true"))
	if err != nil {
		errors <- err
		wg.Done()
		return
	}
	ch <- prometheus.MustNewConstMetric(auc.emailVerified, prometheus.GaugeValue, float64(emailVerified))
	wg.Done()
}

// Describe implements Collector.
func (auc *Auth0UsersCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- auc.totalUsers
}

func (auc *Auth0UsersCollector) Reload(ctx context.Context) error {
	return nil
}

func (auc *Auth0UsersCollector) Update(ctx context.Context, ch chan<- prometheus.Metric) error {
	metricCount := 3
	wg := sync.WaitGroup{}
	wg.Add(metricCount)
	errors := make(chan error, metricCount)
	go auc.updateTotalUsers(&wg, errors, ch)
	go auc.updateBlockedUsers(&wg, errors, ch)
	go auc.updateEmailVerified(&wg, errors, ch)
	wg.Wait()
	close(errors)
	for error := range errors {
		return error
	}
	return nil
}
