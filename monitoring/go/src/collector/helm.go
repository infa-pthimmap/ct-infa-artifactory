package collector

import (
	"fmt"

	"github.com/infa-pthimmap/ct-infa-artifactory/monitoring/go/src/artifactory"
	"github.com/prometheus/client_golang/prometheus"
)

func (e *InfaExporter) ExportHelmStats(ch chan<- prometheus.Metric) error {

	fmt.Println("Exporting Helm Stats")

	helmStats, err := artifactory.VerifyHelmArtDownloads()

	if err != nil {
		return nil
	}

	for metricName, metric := range helmMetrics {
		switch metricName {
		case "helm_status":
			for i := 0; i < len(helmStats); i++ {
				print()
				ch <- prometheus.MustNewConstMetric(metric, prometheus.GaugeValue, float64(helmStats[i].Status), helmStats[i].Slug)
			}

		}
	}
	return nil
}
