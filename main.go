package main

import (
	"flag"
	"fmt"
	"github.com/opensrcit/ftl_exporter/collector"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Version = "development"

var (
	listenAddress string
	metricsPath string
	socket        string
)

func init() {
	flag.StringVar(
		&listenAddress,
		"web.listen-address",
		":9311",
		"Address on which to expose metrics and web interface.")
	flag.StringVar(
		&metricsPath,
		"web.telemetry-path",
		"/metrics",
		"Address on which to expose metrics and web interface.")
	flag.StringVar(&socket, "socket", "/var/run/pihole/FTL.sock", "FTL socket path")

	flag.Usage = func() {
		fmt.Println("FTL Exporter", Version)
		flag.PrintDefaults()
	}

	flag.Parse()
}

func main() {
	log.Println("FTL Exporter", Version)

	ftlExporter := collector.NewFTLExporter(socket)
	prometheus.MustRegister(ftlExporter)

	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`<html lang="en">
             <head><title>FTL Exporter</title></head>
             <body>
             <h1>FTL Exporter</h1>
             <p><a href='/metrics'>Metrics</a></p>
             </body>
             </html>`))
		if err != nil {
			log.Fatal(err)
		}
	})
	log.Println("Listening on", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
