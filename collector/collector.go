package collector

import (
	"context"
	"sync"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/auth0/go-auth0/management"
)

const (
	namespace = "auth0"
)

var (
	factories = make(map[string]func(context.Context, CollectorConfig) (Collector, error))
)

type CollectorConfig struct {
	Logger log.Logger
	Domain *string
	Auth0  *management.Management
}

// Collector is the interface a collector has to implement.
type Collector interface {
	Describe(ch chan<- *prometheus.Desc)
	// Reload configuration, repositories, workflows, etc.
	Reload(context.Context) error
	// Get new metrics and expose them via prometheus registry.
	Update(context.Context, chan<- prometheus.Metric) error
}

type Auth0Collector struct {
	collectors map[string]Collector
	logger     log.Logger
}

type collectorEntry struct {
	name  string
	value Collector
}

func registerCollector(name string, factory func(context.Context, CollectorConfig) (Collector, error)) {
	factories[name] = factory
}

func NewAuth0Collector(config CollectorConfig) (prometheus.Collector, error) {
	wg := sync.WaitGroup{}
	ctx := context.TODO()
	wg.Add(len(factories))
	collectors := make(chan collectorEntry, len(factories))
	errors := make(chan error, len(factories))
	for name, factory := range factories {
		go func(name string, factory func(context.Context, CollectorConfig) (Collector, error)) {
			collector, err := factory(ctx, config)
			if err != nil {
				errors <- err
				wg.Done()
				return
			}
			collectors <- collectorEntry{name, collector}
			wg.Done()
		}(name, factory)
	}
	wg.Wait()
	close(collectors)
	close(errors)
	for error := range errors {
		return nil, error
	}
	collectorMap := make(map[string]Collector)
	for entry := range collectors {
		collectorMap[entry.name] = entry.value
	}
	return &Auth0Collector{
		collectors: collectorMap,
		logger:     config.Logger,
	}, nil
}

// Describe implements the prometheus.Collector interface.
func (ghb Auth0Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, c := range ghb.collectors {
		c.Describe(ch)
	}
}

// Collect implements the prometheus.Collector interface.
func (ghb Auth0Collector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	ctx := context.TODO()
	wg.Add(len(ghb.collectors))
	for name, c := range ghb.collectors {
		go func(name string, c Collector) {
			execute(ctx, name, c, ch, ghb.logger)
			wg.Done()
		}(name, c)
	}
	wg.Wait()
}

func execute(ctx context.Context, name string, c Collector, ch chan<- prometheus.Metric, logger log.Logger) {
	begin := time.Now()
	err := c.Update(ctx, ch)
	duration := time.Since(begin)

	if err != nil {
		_ = level.Error(logger).Log("msg", "collector failed", "name", name, "duration_seconds", duration.Seconds(), "err", err)
	} else {
		_ = level.Debug(logger).Log("msg", "collector succeeded", "name", name, "duration_seconds", duration.Seconds())
	}
}
