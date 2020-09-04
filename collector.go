package main

import (
	"github.com/prometheus/client_golang/prometheus"
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

func init() {
	//Register metrics with prometheus
	prometheus.MustRegister(modelCpuUsageMetric)
	prometheus.MustRegister(modelMemoryUsageMetric)
	prometheus.MustRegister(modelContainersMetric)

	modelCpuUsageMetric.WithLabelValues("seldon-ns","SeldonDeployment","iris").Set(0.019797927555476498)

	modelMemoryUsageMetric.WithLabelValues("seldon-ns","SeldonDeployment","iris").Set(138986.313)

	modelContainersMetric.WithLabelValues("seldon-ns","SeldonDeployment","iris").Set(2.1)

	//TODO: need to set metrics periodically
	//and reset on each iteraction with modelContainersMetric.Reset()
}