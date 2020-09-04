package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/seldonIO/seldon-prometheus-exporter/query"
	"log"
	"time"
)

//Define the metrics we wish to expose
var modelCpuUsageMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "model_cpu_usage_seconds_total",
		Help: "cpu usage for an ML deployment",
	},
	[]string{"namespace","type","name"},
)

var modelMemoryUsageMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "model_memory_usage_bytes",
		Help: "memory usage for an ML deployment",
	},
	[]string{"namespace","type","name"},
)

var modelContainersMetric = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "model_containers_average",
		Help: "Average running containers for an ML deployment",
	},
	[]string{"namespace","type","name"},
)

//TODO: add req and lim

func init() {
	//Register metrics with prometheus
	prometheus.MustRegister(modelCpuUsageMetric)
	prometheus.MustRegister(modelMemoryUsageMetric)
	prometheus.MustRegister(modelContainersMetric)

	// periodically collect metrics
	go func() {
		for {
			collectMetricsUsingQueryTemplate(modelCpuUsageMetric,query.CpuUsageSumSeldonTemplate)
			collectMetricsUsingQueryTemplate(modelMemoryUsageMetric,query.MemUsageSumSeldonTemplate)
			collectMetricsUsingQueryTemplate(modelContainersMetric,query.ContainersUsageSumSeldonTemplate)

			//FIXME: should set this and the Range from an env var, will have to parse string in prom format
			time.Sleep(time.Duration(2 * time.Minute))
		}
	}()
}

func collectMetricsUsingQueryTemplate(metric *prometheus.GaugeVec, queryTemplate string) {
	queryReturnValues, err := query.ObtainMetricValues(queryTemplate, query.DefaultInputData)
	if err != nil {
		log.Println(err)
	}
	metric.Reset()
	for _, v := range queryReturnValues {
		metric.WithLabelValues(v.Namespace, v.ModelType, v.Model).Set(v.Value)
	}
}