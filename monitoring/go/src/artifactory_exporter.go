package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log/level"
	"github.com/infa-pthimmap/ct-infa-artifactory/monitoring/go/src/collector"
	"github.com/infa-pthimmap/ct-infa-artifactory/monitoring/go/src/config"
	"github.com/infa-pthimmap/ct-infa-artifactory/monitoring/go/src/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Errorf("Error creating the config. err: %s", err)
		os.Exit(1)
	}

	exporter, err := collector.NewExporter(conf)
	if err != nil {
		level.Error(conf.Logger).Log("msg", "Error creating an exporter", "err", err)
		os.Exit(1)
	}

	infaexporter, err := collector.InfaNewExporter(conf)
	if err != nil {
		level.Error(conf.Logger).Log("msg", "Error creating an exporter", "err", err)
		os.Exit(1)
	}

	collector.Registry.MustRegister(exporter)
	collector.InfaRegistry.MustRegister(infaexporter)

	level.Info(conf.Logger).Log("msg", "Starting artifactory_exporter", "version", version.Info())

	level.Info(conf.Logger).Log("msg", "Build context", "context", version.BuildContext())

	level.Info(conf.Logger).Log("msg", "Listening on address", "address", conf.ListenAddress)

	//http.Handle(conf.MetricsPath, promhttp.Handler())
	//http.Handle(conf.InfaMetricsPath, infa_prom_http.Handler())

	http.Handle(conf.MetricsPath, promhttp.HandlerFor(collector.Registry, promhttp.HandlerOpts{}))
	http.Handle(conf.InfaMetricsPath, promhttp.HandlerFor(collector.InfaRegistry, promhttp.HandlerOpts{}))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>JFrog Artifactory Exporter</title></head>
             <body>
             <h1>JFrog Exporter</h1>
             <p><a href='` + conf.MetricsPath + `'>Jfrog Metrics</a></p>
			 <p><a href='` + conf.InfaMetricsPath + `'>Infa Metrics</a></p>
             </body>
             </html>`))
	})
	http.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/-/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		metadata := services.ProcessJfrogStatus()

		jsonData, err := json.Marshal(metadata)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the Content-Type header to application/json.
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON data to the response body.
		w.Write(jsonData)
		w.WriteHeader(http.StatusOK)
		//fmt.Fprintf(w, "Status")

	})
	if err := http.ListenAndServe(conf.ListenAddress, nil); err != nil {
		level.Error(conf.Logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}

}
