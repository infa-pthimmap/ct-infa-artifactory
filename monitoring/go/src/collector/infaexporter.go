package collector

import (
	"sync"

	"github.com/go-kit/kit/log"
	"github.com/infa-pthimmap/ct-infa-artifactory/monitoring/go/src/config"
	"github.com/prometheus/client_golang/prometheus"
)

// Exporter collects JFrog Artifactory stats from the given URI and
// exports them using the prometheus metrics package.
type InfaExporter struct {
	artifact_health prometheus.Gauge
	mutex           sync.RWMutex
	logger          log.Logger
}

// NewExporter returns an initialized Exporter.
func InfaNewExporter(conf *config.Config) (*InfaExporter, error) {
	return &InfaExporter{
		artifact_health: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "infa_artifactory_health",
			Help:      "Was the last scrape of artifactory successful.",
		}),
	}, nil
}
