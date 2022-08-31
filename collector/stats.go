package collector

import (
	"context"
	"errors"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/auth0/go-auth0/management"
)

var (
	statsSubsystem = "stats"
)

type Auth0StatsCollector struct {
	config      CollectorConfig
	activeUsers *prometheus.Desc
	// Metrics for the last 24h
	lastDayLogins          *prometheus.Desc
	lastDaySignups         *prometheus.Desc
	lastDayLeakedPasswords *prometheus.Desc
	// Metrics for the last 7 days
	last7dLogins          *prometheus.Desc
	last7dSignups         *prometheus.Desc
	last7dLeakedPasswords *prometheus.Desc
	// Metrics for the last 30d days
	last30dLogins          *prometheus.Desc
	last30dSignups         *prometheus.Desc
	last30dLeakedPasswords *prometheus.Desc
}

func init() {
	registerCollector(statsSubsystem, NewAuth0StatsCollector)
}

// NewOrgActionsCollector returns a new Collector exposing actions billing stats.
func NewAuth0StatsCollector(config CollectorConfig, ctx context.Context) (Collector, error) {
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

	collector := &Auth0StatsCollector{
		config: config,
		activeUsers: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, statsSubsystem, "active_users"),
			"The number of active users that logged in during the last 30 days.",
			nil, prometheus.Labels{"domain": *config.Domain, "tenant": *tenant.FriendlyName},
		),
		lastDayLogins: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, statsSubsystem, "logins"),
			"Number of logins in the interval. The interval is given in days.",
			nil, prometheus.Labels{"domain": *config.Domain, "tenant": *tenant.FriendlyName, "interval": "1"},
		),
		lastDaySignups: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, statsSubsystem, "signups"),
			"Number of signups in the interval. The interval is given in days.",
			nil, prometheus.Labels{"domain": *config.Domain, "tenant": *tenant.FriendlyName, "interval": "1"},
		),
		lastDayLeakedPasswords: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, statsSubsystem, "leaked_passwords"),
			"Number of breached-password detections in the interval. The interval is given in days.",
			nil, prometheus.Labels{"domain": *config.Domain, "tenant": *tenant.FriendlyName, "interval": "1"},
		),
		last7dLogins: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, statsSubsystem, "logins"),
			"Number of logins in the interval. The interval is given in days.",
			nil, prometheus.Labels{"domain": *config.Domain, "tenant": *tenant.FriendlyName, "interval": "7"},
		),
		last7dSignups: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, statsSubsystem, "signups"),
			"Number of signups in the interval. The interval is given in days.",
			nil, prometheus.Labels{"domain": *config.Domain, "tenant": *tenant.FriendlyName, "interval": "7"},
		),
		last7dLeakedPasswords: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, statsSubsystem, "leaked_passwords"),
			"Number of breached-password detections in the interval. The interval is given in days.",
			nil, prometheus.Labels{"domain": *config.Domain, "tenant": *tenant.FriendlyName, "interval": "7"},
		),
		last30dLogins: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, statsSubsystem, "logins"),
			"Number of logins in the interval. The interval is given in days.",
			nil, prometheus.Labels{"domain": *config.Domain, "tenant": *tenant.FriendlyName, "interval": "30"},
		),
		last30dSignups: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, statsSubsystem, "signups"),
			"Number of signups in the interval. The interval is given in days.",
			nil, prometheus.Labels{"domain": *config.Domain, "tenant": *tenant.FriendlyName, "interval": "30"},
		),
		last30dLeakedPasswords: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, statsSubsystem, "leaked_passwords"),
			"Number of breached-password detections in the interval. The interval is given in days.",
			nil, prometheus.Labels{"domain": *config.Domain, "tenant": *tenant.FriendlyName, "interval": "30"},
		),
	}
	err = collector.Reload(ctx)
	if err != nil {
		return nil, err
	}
	return collector, nil
}

// Describe implements Collector.
func (asc *Auth0StatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- asc.activeUsers
	ch <- asc.lastDayLogins
	ch <- asc.lastDaySignups
	ch <- asc.lastDayLeakedPasswords
	ch <- asc.last7dLogins
	ch <- asc.last7dSignups
	ch <- asc.last7dLeakedPasswords
	ch <- asc.last30dLogins
	ch <- asc.last30dSignups
	ch <- asc.last30dLeakedPasswords
}

func (asc *Auth0StatsCollector) Reload(ctx context.Context) error {
	return nil
}

func (asc *Auth0StatsCollector) Update(ctx context.Context, ch chan<- prometheus.Metric) error {
	activeUsers, err := asc.config.Auth0.Stat.ActiveUsers()
	if err != nil {
		return err
	}
	now := time.Now().Truncate(time.Duration(time.Hour * 24))
	from := now.AddDate(0, 0, -30).Format("20060102")
	to := now.Format("20060102")
	stats, err := asc.config.Auth0.Stat.Daily(
		management.Parameter("from", from),
		management.Parameter("to", to),
	)
	if err != nil {
		return err
	}
	lastDayLogins := 0
	lastDaySignups := 0
	lastDayLeakedPasswords := 0
	last7dLogins := 0
	last7dSignups := 0
	last7dLeakedPasswords := 0
	last30dLogins := 0
	last30dSignups := 0
	last30dLeakedPasswords := 0

	for index, stat := range stats {
		if (index + 1) == len(stats) {
			lastDayLogins = *stat.Logins
			lastDaySignups = *stat.Signups
			lastDayLeakedPasswords = *stat.LeakedPasswords
		}
		if (index + 7) >= len(stats) {
			last7dLogins += *stat.Logins
			last7dSignups += *stat.Signups
			last7dLeakedPasswords += *stat.LeakedPasswords
		}
		last30dLogins += *stat.Logins
		last30dSignups += *stat.Signups
		last30dLeakedPasswords += *stat.LeakedPasswords
	}

	ch <- prometheus.MustNewConstMetric(asc.activeUsers, prometheus.GaugeValue, float64(activeUsers))
	ch <- prometheus.MustNewConstMetric(asc.lastDayLogins, prometheus.GaugeValue, float64(lastDayLogins))
	ch <- prometheus.MustNewConstMetric(asc.lastDaySignups, prometheus.GaugeValue, float64(lastDaySignups))
	ch <- prometheus.MustNewConstMetric(asc.lastDayLeakedPasswords, prometheus.GaugeValue, float64(lastDayLeakedPasswords))
	ch <- prometheus.MustNewConstMetric(asc.last7dLogins, prometheus.GaugeValue, float64(last7dLogins))
	ch <- prometheus.MustNewConstMetric(asc.last7dSignups, prometheus.GaugeValue, float64(last7dSignups))
	ch <- prometheus.MustNewConstMetric(asc.last7dLeakedPasswords, prometheus.GaugeValue, float64(last7dLeakedPasswords))
	ch <- prometheus.MustNewConstMetric(asc.last30dLogins, prometheus.GaugeValue, float64(last30dLogins))
	ch <- prometheus.MustNewConstMetric(asc.last30dSignups, prometheus.GaugeValue, float64(last30dSignups))
	ch <- prometheus.MustNewConstMetric(asc.last30dLeakedPasswords, prometheus.GaugeValue, float64(last30dLeakedPasswords))

	return nil
}
