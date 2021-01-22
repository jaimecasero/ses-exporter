package main

import (
	"flag"
	"net/http"
	"ses-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

var (
	listenAddress = flag.String("web.listen-address", ":9435",
		"Address to listen on for telemetry")
	metricsPath = flag.String("web.telemetry-path", "/metrics",
		"Path under which to expose metrics")
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infoln(r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()

	prometheus.MustRegister(collector.NewExporter())

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
               <head><title>AWS SES Exporter</title></head>
               <body>
               <h1>AWS SES Exporter</h1>
               <p><a href='` + *metricsPath + `'>Metrics</a></p>
               </body>
               </html>`))
	})
	log.Infoln("Starting HTTP server on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, logRequest(http.DefaultServeMux)))
}
