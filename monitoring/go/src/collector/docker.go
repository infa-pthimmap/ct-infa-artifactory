package collector

import (
	"fmt"

	"github.com/infa-pthimmap/ct-infa-artifactory/monitoring/go/src/artifactory"
	"github.com/prometheus/client_golang/prometheus"
)

func (e *InfaExporter) ExportDockerStats(ch chan<- prometheus.Metric) error {

	fmt.Println("Exporting docker Stats")

	dockerStats, err := artifactory.VerifyDockerDownloads()

	if err != nil {
		return nil
	}

	for metricName, metric := range dockerMetrics {
		switch metricName {
		case "docker_status":
			for i := 0; i < len(dockerStats); i++ {
				ch <- prometheus.MustNewConstMetric(metric, prometheus.GaugeValue, float64(dockerStats[i].Status), dockerStats[i].Slug)
			}
		}
	}
	return nil
}
