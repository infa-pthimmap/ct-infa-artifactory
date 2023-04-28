package collector

import (
	"fmt"

	"github.com/infa-pthimmap/ct-infa-artifactory/monitoring/go/src/artifactory"
	"github.com/prometheus/client_golang/prometheus"
)

func (e *InfaExporter) ExportCertificateStats(ch chan<- prometheus.Metric) error {

	fmt.Println("Exporting Certificate Stats")

	certStats, err := artifactory.GetCertificatesDetails()

	if err != nil {
		return nil
	}

	for metricName, metric := range certMetrics {
		switch metricName {
		case "cert_status":
			ch <- prometheus.MustNewConstMetric(metric, prometheus.GaugeValue, float64(certStats.DaysToExpire), fmt.Sprintf("%d", certStats.DaysToExpire), certStats.ExpiresOn)
		case "artifactory_status":
			ch <- prometheus.MustNewConstMetric(metric, prometheus.GaugeValue, float64(certStats.Status), certStats.ArtifcatoryUrl)

		}
	}
	return nil
}
