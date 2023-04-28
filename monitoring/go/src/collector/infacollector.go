package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/version"
)

const (
	ns = "infa_artifactory"
)

func newInfaMetric(metricName string, subsystem string, docString string, labelNames []string) *prometheus.Desc {
	return prometheus.NewDesc(prometheus.BuildFQName(ns, subsystem, metricName), docString, labelNames, nil)
}

type infaMetrics map[string]*prometheus.Desc

var (
	infaDefaultLabelNames = []string{"repo_name"}
	infaHealthLabelNames  = []string{"artifactory_url"}
	infaCertLabelNames    = []string{"days_to_expire", "expires_on"}
	dockerMetrics         = infaMetrics{
		"docker_status": newInfaMetric("status", "docker", "docker stats.", infaDefaultLabelNames),
		//"helm_status":   newInfaMetric("status", "helm", "helm stats", append([]string{"status"}, infaDefaultLabelNames...)),
	}

	helmMetrics = infaMetrics{
		"helm_status": newInfaMetric("status", "helm", "helm stats.", infaDefaultLabelNames),
	}

	certMetrics = infaMetrics{
		"cert_status":        newInfaMetric("status", "certificate", "certificate stats", infaCertLabelNames),
		"artifactory_status": newInfaMetric("status", "health", "certificate stats", infaHealthLabelNames),
	}

	InfaRegistry = prometheus.NewRegistry()
)

func init() {
	//prometheus.MustRegister(version.NewCollector("infa_artifactory_exporter"))
	InfaRegistry.MustRegister(version.NewCollector("infa_artifactory_exporter"))
}

// Describe describes all the metrics ever exported by the Artifactory exporter. It
// implements prometheus.Collector.
func (e *InfaExporter) Describe(ch chan<- *prometheus.Desc) {

	for _, m := range dockerMetrics {
		ch <- m
	}

	for _, m := range helmMetrics {
		ch <- m
	}

	for _, m := range certMetrics {
		ch <- m
	}

	ch <- e.artifact_health.Desc()
}

// Collect fetches the stats from  Artifactory and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *InfaExporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	artifact_health := e.scrape(ch)
	ch <- e.artifact_health
	e.artifact_health.Set(artifact_health)
}

func (e *InfaExporter) scrape(ch chan<- prometheus.Metric) (up float64) {

	docker_err := e.ExportDockerStats(ch)
	if docker_err != nil {
		return 0
	}

	helm_err := e.ExportHelmStats(ch)
	if helm_err != nil {
		return 0
	}

	cert_err := e.ExportCertificateStats(ch)
	if cert_err != nil {
		return 0
	}

	return 1
}
