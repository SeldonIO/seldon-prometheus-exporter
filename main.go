package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	//at some point could do https://stackoverflow.com/a/50745945/9705485 to only see state metrics
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}