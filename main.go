package main

import (
	"fmt"
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
	http.Handle("/healthz", handle(HealthzHandler))
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type handle func(w http.ResponseWriter, req *http.Request) error

type Error struct {
	Code    int
	Message string
}

func (e Error) Error() string {
	if e.Message == "" {
		e.Message = http.StatusText(e.Code)
	}
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func (h handle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Handler panic: %v", r)
		}
	}()
	if err := h(w, req); err != nil {
		log.Printf("Handler error: %v", err)

		if httpErr, ok := err.(Error); ok {
			http.Error(w, httpErr.Message, httpErr.Code)
		}
	}
}

func HealthzHandler(w http.ResponseWriter, req *http.Request) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("all is well"))
	return nil
}